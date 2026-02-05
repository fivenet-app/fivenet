import type { UseWebSocketReturn } from '@vueuse/core';
import { writeUInt32BE } from '~/utils/array';
import { Body, Cancel, GrpcFrame, Header, HeaderValue } from '~~/gen/ts/resources/grpcws/grpcws';
import { headersToMetadata } from '../../bridge/utils';
import { errCancelled, errInternal } from '../../errors';
import type { Metadata } from '../../metadata';
import type { Transport, TransportFactory, TransportOptions } from '../transport';
import { createRpcError } from './utils';

const websocketChannelMaxStreamCount = 7;

// Control stream is reserved and never used for RPC calls.
const CONTROL_STREAM_ID = 0 as const;

// Control-plane "operations" sent in Header.operation on stream 0.
const CONTROL_OP_AUTH = 'auth' as const;
const CONTROL_OP_REAUTH = 'reauth' as const;
const CONTROL_OP_AUTH_OK = 'auth_ok' as const;

type TokenProvider = () => string | null;

type AuthState =
    | { kind: 'none' }
    | { kind: 'pending'; promise: Promise<void>; resolve: () => void; reject: (err: Error) => void }
    | { kind: 'ok' }
    | { kind: 'failed'; err: Error };

// Default token provider: change the key to your real sessionStorage key.
function defaultTokenProvider(): string | null {
    const raw = sessionStorage.getItem('fivenet:user_token_v1');
    if (!raw) return null;
    return raw;
}

export function WebsocketChannelTransport(
    logger: ILogger,
    webSocket: UseWebSocketReturn<ArrayBuffer>,
    tokenProvider: TokenProvider = defaultTokenProvider,
): TransportFactory {
    const wsChannel = new WebsocketChannelImpl(logger, webSocket, tokenProvider);

    return (opts: TransportOptions) => {
        opts.debug && logger.debug('Websocket factory triggered, status:', webSocket.status.value);
        if (webSocket.status.value === 'CLOSED') {
            webSocket.open();
        }

        return wsChannel.getStream(opts);
    };
}

interface GrpcStream extends Transport {
    readonly streamId: number;
    readonly service: string;
    readonly method: string;
    readonly isStream: boolean;

    closed: boolean;
}

interface WebsocketChannel {
    getStream(options: TransportOptions): GrpcStream;
}

class WebsocketChannelImpl implements WebsocketChannel {
    private logger: ILogger;
    protected ws: UseWebSocketReturn<ArrayBuffer>;
    private tokenProvider: TokenProvider;

    readonly activeStreams = new Map<number, [TransportOptions, WebsocketChannelStream]>();

    public lastStreamId = 1;

    private authState: AuthState = { kind: 'none' };

    constructor(logger: ILogger, ws: UseWebSocketReturn<ArrayBuffer>, tokenProvider: TokenProvider) {
        this.logger = logger;
        this.ws = ws;
        this.tokenProvider = tokenProvider;

        watch(ws.data, async (val) => this.onMessage(val));
        watchThrottled(
            ws.status,
            (val) => {
                if (val === 'OPEN') return;

                // Close all streams when the websocket connection is lost/closed
                this.activeStreams.forEach((as) => {
                    as[1].cancel();
                    as[1].closed = true;
                    this.activeStreams.delete(as[1].streamId);
                });

                this.lastStreamId = 1;
                this.authState = { kind: 'none' };
            },
            {
                immediate: true,
                throttle: 250, // 250ms throttle to avoid rapid status changes
            },
        );
    }

    async onMessage(data: ArrayBuffer | null): Promise<void> {
        if (data === null) return;

        const frame = GrpcFrame.fromBinary(new Uint8Array(data));
        const streamId = frame.streamId;

        if (frame.payload.oneofKind === 'ping') {
            this.logger.debug('Received websocket ping, pong:', frame.payload.ping.pong);
            return;
        }

        // CONTROL stream handling (streamId = 0)
        if (streamId === CONTROL_STREAM_ID) {
            this.onControlFrame(frame);
            return;
        }

        const stream = this.activeStreams.get(streamId);
        if (!stream) {
            // If the stream does not exist, we ignore `complete` messages
            if (frame.payload.oneofKind !== 'complete') {
                this.logger.warn('Stream does not exist', streamId, `${frame.payload.oneofKind}`, 'ignoring message');
            }
            return;
        }

        const wsStream: WebsocketChannelStream = stream[1];
        switch (frame.payload.oneofKind) {
            case 'header': {
                stream[0].debug &&
                    this.logger.debug('Received header for stream', streamId, `${stream[1].service}/${stream[1].method}`);

                const header = frame.payload.header;
                if (header === null) return;

                const metaData = headersToMetadata(header.headers);
                stream[0].onHeaders(metaData, header.status);

                const err = createRpcError(metaData, stream[0].methodDefinition);
                if (err) {
                    stream[0].onEnd(err);
                }
                break;
            }

            case 'body': {
                stream[0].debug &&
                    this.logger.debug('Received body for stream', streamId, `${stream[1].service}/${stream[1].method}`);

                const body = frame.payload.body;
                if (body === null) return;

                stream[0].onChunk(body.data);
                break;
            }

            case 'complete': {
                stream[0].debug &&
                    this.logger.debug('Received complete for stream', streamId, `${stream[1].service}/${stream[1].method}`);

                stream[0].onEnd();
                wsStream.closed = true;
                // Remove completed stream
                this.activeStreams.delete(streamId);
                break;
            }

            case 'failure': {
                const failure = frame.payload.failure;
                if (failure === null) {
                    stream[0].onEnd(errInternal);
                    return;
                }

                stream[0].debug &&
                    this.logger.debug(
                        'Received failure for stream',
                        streamId,
                        `${stream[1].service}/${stream[1].method}`,
                        'status:',
                        failure.errorStatus,
                        'msg:',
                        failure.errorMessage,
                    );

                const metaData = headersToMetadata(failure.headers);
                const err = createRpcError(metaData, stream[0].methodDefinition);
                stream[0].onEnd(err);
                this.activeStreams.delete(streamId);
                break;
            }

            case 'cancel': {
                stream[0].debug &&
                    this.logger.debug('Received cancel for stream', streamId, `${stream[1].service}/${stream[1].method}`);

                stream[0].onEnd(errCancelled);
                break;
            }

            default:
                stream[0].debug &&
                    this.logger.debug(
                        'Received unknown message type for stream',
                        streamId,
                        frame.payload.oneofKind,
                        `${stream[1].service}/${stream[1].method}`,
                    );
                break;
        }
    }

    private onControlFrame(frame: GrpcFrame) {
        // Control frames use Header + Failure
        switch (frame.payload.oneofKind) {
            case 'header': {
                const h = frame.payload.header;
                if (!h) return;

                if (h.operation === CONTROL_OP_AUTH_OK) {
                    this.logger.debug('WS auth operation: ok');
                    if (this.authState.kind === 'pending') this.authState.resolve();
                    this.authState = { kind: 'ok' };
                    return;
                }

                // If server responds with something else on control stream, treat as failure.
                const err = new Error(`WS control header unexpected operation: ${h.operation}`);
                if (this.authState.kind === 'pending') this.authState.reject(err);
                this.authState = { kind: 'failed', err };
                return;
            }

            case 'failure': {
                const f = frame.payload.failure;
                const msg = f?.errorMessage || 'Unauthenticated';
                const err = new Error(`WS auth failed: ${msg}`);
                this.logger.warn(err.message);

                if (this.authState.kind === 'pending') this.authState.reject(err);
                this.authState = { kind: 'failed', err };
                return;
            }

            default:
                // Ignore others (complete/cancel/body shouldn't happen on control stream)
                return;
        }
    }

    async sendToWebsocket(opts: TransportOptions, toSend: GrpcFrame, usingBuffer: boolean = true): Promise<boolean> {
        if (toSend.streamId !== CONTROL_STREAM_ID && !this.activeStreams.has(toSend.streamId)) {
            opts.debug && this.logger.debug('sendToWs: Stream does not exist', toSend.streamId);
            return false;
        }

        return this.ws.send(GrpcFrame.toBinary(toSend).buffer as ArrayBuffer, usingBuffer);
    }

    private async sendControlHeader(operation: string, headers: Record<string, string[]>) {
        const header = Header.create();
        header.operation = operation;

        const headerMap = header.headers;
        for (const [key, values] of Object.entries(headers)) {
            const hv = HeaderValue.create();
            hv.value = values;
            headerMap[key] = hv;
        }

        // Control sends don't belong to a normal stream in activeStreams,
        // so use a dummy TransportOptions for logging only.
        this.ws.send(
            GrpcFrame.toBinary(
                GrpcFrame.create({
                    streamId: CONTROL_STREAM_ID,
                    payload: { oneofKind: 'header', header },
                }),
            ).buffer as ArrayBuffer,
            true,
        );
    }

    async ensureAuthenticated(): Promise<void> {
        if (this.ws.status.value !== 'OPEN')
            // Let caller open first; once open, auth will be attempted.
            throw new Error('WebSocket not open');

        if (this.authState.kind === 'ok') return;
        if (this.authState.kind === 'failed') throw this.authState.err;
        if (this.authState.kind === 'pending') return this.authState.promise;

        const token = this.tokenProvider();
        if (!token) throw new Error('Missing character token (sessionStorage)');

        let resolve!: () => void;
        let reject!: (err: Error) => void;
        const promise = new Promise<void>((res, rej) => {
            resolve = res;
            reject = rej;
        });

        this.authState = { kind: 'pending', promise, resolve, reject };

        await this.sendControlHeader(CONTROL_OP_AUTH, {
            Authorization: [`Bearer ${token}`],
        });

        return promise;
    }

    async reauth(): Promise<void> {
        if (this.ws.status.value !== 'OPEN') throw new Error('WebSocket not open');

        const token = this.tokenProvider();
        if (!token) throw new Error('Missing character token (sessionStorage)');

        // If you want to force reauth to await auth_ok, set authState pending here too.
        await this.sendControlHeader(CONTROL_OP_REAUTH, {
            Authorization: [`Bearer ${token}`],
        });
    }

    public getNextStreamId(): number {
        // Reset stream id back to 1 if max is reached
        if (this.lastStreamId >= websocketChannelMaxStreamCount) {
            this.lastStreamId = 1;
            return this.lastStreamId;
        }

        return this.lastStreamId++;
    }

    getStream(opts: TransportOptions): GrpcStream {
        return new WebsocketChannelStream(this, this.logger, this.getNextStreamId(), opts);
    }
}

/**
 * WebsocketChannelStream: gRPC-WebSocket stream.
 *
 * Usage:
 *   - Higher-level code should handle resending any initial messages if needed.
 */
class WebsocketChannelStream {
    wsChannel: WebsocketChannelImpl;
    logger: ILogger;
    streamId: number;
    opts: TransportOptions;
    service: string;
    method: string;
    isStream: boolean;

    closed: boolean = false;

    constructor(wsChannel: WebsocketChannelImpl, logger: ILogger, streamId: number, opts: TransportOptions) {
        this.wsChannel = wsChannel;
        this.logger = logger;
        this.streamId = streamId;
        this.opts = opts;
        this.service = opts.methodDefinition.service.typeName;
        this.method = opts.methodDefinition.name;
        this.isStream = opts.methodDefinition.serverStreaming || opts.methodDefinition.clientStreaming;
    }

    async start(metadata: Metadata) {
        this.opts.debug && this.logger.debug('Stream start', this.streamId, `${this.service}/${this.method}`);

        // Ensure control-plane auth succeeded before creating a real stream.
        this.wsChannel.ensureAuthenticated();

        this.wsChannel.activeStreams.set(this.streamId, [this.opts, this]);

        const header = Header.create();
        /**
         * Header.operation multiplexing:
         * - For RPC streams (streamId > 0): operation is "Service/Method"
         * - For control stream (streamId == 0): operation is one of "auth", "reauth", "auth_ok"
         */
        header.operation = `${this.service}/${this.method}`;

        const headerMap = header.headers;
        // Only accept Metadata for metadata
        metadata.forEach((key, values) => {
            const headerValue = HeaderValue.create();
            headerValue.value = values;
            headerMap[key] = headerValue;
        });

        this.wsChannel.sendToWebsocket(
            this.opts,
            GrpcFrame.create({
                streamId: this.streamId,
                payload: {
                    oneofKind: 'header',
                    header: header,
                },
            }),
        );
    }

    async sendMessage(msgBytes: Uint8Array, complete?: boolean) {
        this.opts.debug && this.logger.debug('Stream send', this.streamId, `${this.service}/${this.method}`);

        const output = new Uint8Array(msgBytes.length + 5);
        output[0] = 0; // Compression none
        writeUInt32BE(output, msgBytes.length, 1);
        output.set(msgBytes, 5);

        const body = Body.create({
            data: output,
            complete: !!complete,
        });

        await this.wsChannel.sendToWebsocket(
            this.opts,
            GrpcFrame.create({
                streamId: this.streamId,
                payload: {
                    oneofKind: 'body',
                    body: body,
                },
            }),
        );
    }

    async finishSend() {
        this.opts.debug && this.logger.debug('Stream complete', this.streamId, `${this.service}/${this.method}`);

        await this.wsChannel.sendToWebsocket(
            this.opts,
            GrpcFrame.create({
                streamId: this.streamId,
                payload: {
                    oneofKind: 'complete',
                    complete: {},
                },
            }),
        );
    }

    async cancel(err?: Error) {
        if (this.closed) return;

        this.opts.debug && this.logger.debug('Stream cancel', this.streamId, `${this.service}/${this.method}`, 'Error:', err);

        this.opts.onEnd(err ?? errCancelled);

        await this.wsChannel.sendToWebsocket(
            this.opts,
            GrpcFrame.create({
                streamId: this.streamId,
                payload: {
                    oneofKind: 'cancel',
                    cancel: Cancel.create(),
                },
            }),
        );
    }
}

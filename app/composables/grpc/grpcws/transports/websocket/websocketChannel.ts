import type { UseWebSocketReturn } from '@vueuse/core';
import { writeUInt32BE } from '~/utils/array';
import { Body, Cancel, GrpcFrame, Header, HeaderValue } from '~~/gen/ts/resources/common/grpcws/grpcws';
import { headersToMetadata } from '../../bridge/utils';
import { errCancelled, errInternal } from '../../errors';
import type { Metadata } from '../../metadata';
import type { Transport, TransportFactory, TransportOptions } from '../transport';
import { createRpcError } from './utils';

const websocketChannelMaxStreamCount = 7;

export function WebsocketChannelTransport(logger: ILogger, webSocket: UseWebSocketReturn<ArrayBuffer>): TransportFactory {
    const wsChannel = new WebsocketChannelImpl(logger, webSocket);

    return (opts: TransportOptions) => {
        opts.debug && logger.debug('Websocket factory triggered');
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
    close(): void;
}

class WebsocketChannelImpl implements WebsocketChannel {
    private logger: ILogger;
    protected ws: UseWebSocketReturn<ArrayBuffer>;
    readonly activeStreams = new Map<number, [TransportOptions, GrpcStream]>();
    protected lastStreamId = 0;

    constructor(logger: ILogger, ws: UseWebSocketReturn<ArrayBuffer>) {
        this.logger = logger;
        this.ws = ws;
        watch(ws.data, async (val) => this.onMessage(val));
        watchThrottled(
            ws.status,
            () => {
                if (ws.status.value === 'CLOSED') {
                    this.activeStreams.forEach((as) => {
                        as[1].cancel();
                        as[1].closed = true;

                        this.activeStreams.delete(as[1].streamId);
                    });
                }
            },
            {
                throttle: 500,
            },
        );
    }

    close() {
        this.ws.close();
    }

    async onMessage(data: ArrayBuffer | null): Promise<void> {
        if (data === null) {
            return;
        }

        const frame = GrpcFrame.fromBinary(new Uint8Array(data));
        const streamId = frame.streamId;
        if (frame.payload.oneofKind === 'ping') {
            this.logger.debug('Received websocket ping, pong:', frame.payload.ping.pong);
            return;
        }

        const stream = this.activeStreams.get(streamId);

        if (stream) {
            switch (frame.payload.oneofKind) {
                case 'header': {
                    stream[0].debug &&
                        this.logger.debug('Received header for stream', streamId, `${stream[1].service}/${stream[1].method}`);

                    const header = frame.payload.header;
                    if (header === null) {
                        return;
                    }

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
                    if (body === null) {
                        return;
                    }

                    stream[0].onChunk(body.data);
                    break;
                }

                case 'complete': {
                    stream[0].debug &&
                        this.logger.debug('Received complete for stream', streamId, `${stream[1].service}/${stream[1].method}`);

                    stream[0].onEnd();
                    stream[1].closed = true;
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
                    stream[0].onEnd(createRpcError(metaData, stream[0].methodDefinition));
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
        } else {
            this.logger.warn('Stream does not exist', streamId);
        }
    }

    async sendToWebsocket(opts: TransportOptions, toSend: GrpcFrame, usingBuffer: boolean = true): Promise<boolean> {
        if (!this.activeStreams.has(toSend.streamId)) {
            opts.debug && this.logger.debug('sendToWs: Stream does not exist', toSend.streamId);
            return false;
        }

        return this.ws.send(GrpcFrame.toBinary(toSend).buffer as ArrayBuffer, usingBuffer);
    }

    private getNextStreamId(): number {
        // Reset stream id back to 0 if max is reached
        if (this.lastStreamId >= websocketChannelMaxStreamCount) {
            this.lastStreamId = 0;
            return this.lastStreamId;
        }

        return this.lastStreamId++;
    }

    getStream(opts: TransportOptions): GrpcStream {
        return new WebsocketChannelStream(this, this.logger, this.getNextStreamId(), opts);
    }
}

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

    start(metadata: Metadata) {
        this.opts.debug && this.logger.debug('Stream start', this.streamId, `${this.service}/${this.method}`);

        this.wsChannel.activeStreams.set(this.streamId, [this.opts, this]);

        const header = Header.create();
        header.operation = `${this.service}/${this.method}`;
        const headerMap = header.headers;
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

        this.wsChannel.sendToWebsocket(
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

        this.wsChannel.sendToWebsocket(
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
        if (this.closed) {
            return;
        }

        this.opts.debug && this.logger.debug('Stream cancel', this.streamId, `${this.service}/${this.method}`);

        this.opts.onEnd(err ?? errCancelled);

        this.wsChannel.sendToWebsocket(
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

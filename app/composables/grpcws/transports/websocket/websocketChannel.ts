import type { UseWebSocketReturn } from '@vueuse/core';
import { writeUInt32BE } from '~/utils/array';
import { Body, Cancel, Complete, GrpcFrame, Header, HeaderValue } from '~~/gen/ts/resources/common/grpcws/grpcws';
import { headersToMetadata } from '../../bridge/utils';
import { errCancelled, errInternal, errUnavailable } from '../../errors';
import { Metadata } from '../../metadata';
import type { Transport, TransportFactory, TransportOptions } from '../transport';
import { createRpcError } from './utils';

export function WebsocketChannelTransport(logger: ILogger, webSocket: UseWebSocketReturn<ArrayBuffer>): TransportFactory {
    const wsChannel = new WebsocketChannelImpl(logger, webSocket);

    return (opts: TransportOptions) => {
        opts.debug && logger.debug('Websocket factory triggered');
        if (webSocket.status.value === 'CLOSED') {
            webSocket.open();

            // (Re-)start any active stream channels
            if (wsChannel.activeStreams.size > 0) {
                wsChannel.activeStreams.forEach((stream) => {
                    if (!stream[1].isStream) {
                        return;
                    }

                    stream[1].start(new Metadata());
                });
            }
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
    protected lastStreamId = 1;

    constructor(logger: ILogger, ws: UseWebSocketReturn<ArrayBuffer>) {
        this.logger = logger;
        this.ws = ws;
        watch(ws.data, async (val) => this.onMessage(val));
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
                    stream[0].debug && this.logger.debug('Received header for stream', streamId);

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
                    stream[0].debug && this.logger.debug('Received body for stream', streamId);

                    const body = frame.payload.body;
                    if (body === null) {
                        return;
                    }

                    stream[0].onChunk(body.data);
                    break;
                }

                case 'complete': {
                    stream[0].debug && this.logger.debug('Received complete for stream', streamId);

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
                        this.logger.debug('Received failure for stream', streamId, failure.errorStatus, failure.errorMessage);

                    const metaData = headersToMetadata(failure.headers);
                    stream[0].onEnd(createRpcError(metaData, stream[0].methodDefinition));
                    break;
                }

                case 'cancel': {
                    stream[0].debug && this.logger.debug('Received cancel for stream', streamId);

                    stream[0].onEnd(errCancelled);
                    break;
                }

                default:
                    stream[0].debug &&
                        this.logger.debug('Received unknown message type for stream', streamId, frame.payload.oneofKind);
                    break;
            }
        } else {
            this.logger.warn('Stream does not exist', streamId);
        }
    }

    async sendToWebsocket(opts: TransportOptions, toSend: GrpcFrame): Promise<void> {
        if (!this.activeStreams.has(toSend.streamId)) {
            opts.debug && this.logger.debug('Stream does not exist', toSend.streamId);
            return;
        }

        if (this.ws.status.value === 'CLOSED') {
            throw errUnavailable;
        }

        this.ws.send(GrpcFrame.toBinary(toSend), true);
    }

    getStream(opts: TransportOptions): GrpcStream {
        const currentStreamId = this.lastStreamId++;

        return new WebsocketChannelStream(this, this.logger, currentStreamId, opts);
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
        this.opts.debug && this.logger.debug('Stream send', this.streamId);

        const output = new Uint8Array(msgBytes.length + 5);
        output[0] = 0; // Compression none
        writeUInt32BE(output, msgBytes.length, 1);
        output.set(msgBytes, 5);

        const body = Body.create();
        body.data = output;
        body.complete = !!complete;

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
        this.opts.debug && this.logger.debug('Stream complete', this.streamId);

        this.wsChannel.sendToWebsocket(
            this.opts,
            GrpcFrame.create({
                streamId: this.streamId,
                payload: {
                    oneofKind: 'complete',
                    complete: Complete.create(),
                },
            }),
        );
    }

    async cancel() {
        if (this.closed) {
            return;
        }

        this.opts.debug && this.logger.debug('Stream cancel', this.streamId);

        this.opts.onEnd(errCancelled);

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

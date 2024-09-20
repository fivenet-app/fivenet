import type { UseWebSocketReturn } from '@vueuse/core';
import { writeUInt32BE } from '~/utils/array';
import { Body, Cancel, Complete, GrpcFrame, Header, HeaderValue } from '~~/gen/ts/resources/common/grpcws/grpcws';
import { headersToMetadata } from '../../bridge/utils';
import { errCancelled, errInternal, errUnavailable } from '../../errors';
import { Metadata } from '../../metadata';
import { type Transport, type TransportFactory, type TransportOptions } from '../transport';
import { createRpcError } from './utils';

export function WebsocketChannelTransport(logger: ILogger, webSocket: UseWebSocketReturn<any>): TransportFactory {
    const wsChannel = new WebsocketChannelImpl(logger, webSocket);

    return (opts: TransportOptions) => {
        opts.debug && logger.debug('Websocket factory triggered');
        if (webSocket.status.value === 'CLOSED') {
            webSocket.open();

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
}

interface WebsocketChannel {
    getStream(options: TransportOptions): GrpcStream;
    close(): void;
}

class WebsocketChannelImpl implements WebsocketChannel {
    private logger: ILogger;
    protected ws: UseWebSocketReturn<any>;
    readonly activeStreams = new Map<number, [TransportOptions, GrpcStream]>();
    protected lastStreamId = 1;

    constructor(logger: ILogger, ws: UseWebSocketReturn<any>) {
        this.logger = logger;
        this.ws = ws;
        watch(ws.data, async (val) => this.onMessage(val as any));
    }

    async close() {
        this.ws.close();
    }

    async onMessage(data: ArrayBuffer): Promise<void> {
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

    getStream(opts: TransportOptions): GrpcStream {
        const currentStreamId = this.lastStreamId++;
        const self = this;

        async function sendToWebsocket(toSend: GrpcFrame): Promise<void> {
            if (!self.activeStreams.has(toSend.streamId)) {
                opts.debug && self.logger.debug('Stream does not exist', toSend.streamId);
                return;
            }

            if (self.ws.status.value === 'CLOSED') {
                throw errUnavailable;
            }

            self.ws.send(GrpcFrame.toBinary(toSend), true);
        }

        function newFrame(): GrpcFrame {
            const frame = GrpcFrame.create();
            frame.streamId = currentStreamId;
            return frame;
        }

        // Question: can this structure be reused or is it one time use?
        const stream = {
            streamId: currentStreamId,
            service: opts.methodDefinition.service.typeName,
            method: opts.methodDefinition.name,
            isStream: opts.methodDefinition.serverStreaming || opts.methodDefinition.clientStreaming,

            start: (metadata: Metadata) => {
                opts.debug &&
                    this.logger.debug(
                        'Stream start',
                        currentStreamId,
                        `${opts.methodDefinition.service.typeName}/${opts.methodDefinition.name}`,
                    );

                self.activeStreams.set(currentStreamId, [opts, stream]);

                const header = Header.create();
                header.operation = `${opts.methodDefinition.service.typeName}/${opts.methodDefinition.name}`;
                const headerMap = header.headers;
                metadata.forEach((key, values) => {
                    const headerValue = HeaderValue.create();
                    headerValue.value = values;
                    headerMap[key] = headerValue;
                });

                const frame = newFrame();
                frame.payload = {
                    oneofKind: 'header',
                    header: header,
                };

                sendToWebsocket(frame);
            },

            sendMessage: async (msgBytes: Uint8Array, complete?: boolean) => {
                opts.debug && this.logger.debug('Stream send', currentStreamId);

                const output = new Uint8Array(msgBytes.length + 5);
                output[0] = 0; // Compression none
                writeUInt32BE(output, msgBytes.length, 1);
                output.set(msgBytes, 5);

                const body = Body.create();
                body.data = output;
                body.complete = !!complete;

                const frame = newFrame();
                frame.payload = {
                    oneofKind: 'body',
                    body: body,
                };

                sendToWebsocket(frame);
            },

            finishSend: async () => {
                opts.debug && this.logger.debug('Stream complete', currentStreamId);

                const frame = newFrame();
                frame.payload = {
                    oneofKind: 'complete',
                    complete: Complete.create(),
                };

                sendToWebsocket(frame);
            },

            cancel: async () => {
                opts.debug && this.logger.debug('Stream cancel', currentStreamId);

                const frame = newFrame();
                frame.payload = {
                    oneofKind: 'cancel',
                    cancel: Cancel.create(),
                };

                opts.onEnd(errCancelled);

                sendToWebsocket(frame);
            },
        };

        return stream;
    }
}

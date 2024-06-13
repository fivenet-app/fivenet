import type { UseWebSocketReturn } from '@vueuse/core';
import { writeUInt32BE } from '~/utils/array';
import { Body, Cancel, Complete, GrpcFrame, Header, HeaderValue } from '~~/gen/ts/resources/common/grpcws/grpcws';
import { headersToMetadata } from '../../bridge/utils';
import { Metadata } from '../../metadata';
import { type Transport, type TransportFactory, type TransportOptions } from '../transport';
import { createRpcError } from './utils';
import { errCancelled, errInternal, errUnavailable } from '../../errors';

export function WebsocketChannelTransport(webSocket: UseWebSocketReturn<any>): TransportFactory {
    const wsChannel = new WebsocketChannelImpl(webSocket);

    return (opts: TransportOptions) => {
        opts.debug && console.debug('GRPC-WS: Websocket factory triggered');
        if (webSocket.status.value === 'CLOSED') {
            webSocket.open();
        }

        return wsChannel.getStream(opts);
    };
}

interface GrpcStream extends Transport {
    readonly streamId: number;
}

interface WebsocketChannel {
    getStream(options: TransportOptions): GrpcStream;
    close(): void;
}

class WebsocketChannelImpl implements WebsocketChannel {
    protected ws: UseWebSocketReturn<any>;
    readonly activeStreams = new Map<number, [TransportOptions, GrpcStream]>();
    protected streamId = 0;

    constructor(ws: UseWebSocketReturn<any>) {
        this.ws = ws;
        watch(ws.data, async (val) => this.onMessage(val as any));
    }

    async close() {
        this.ws.close();
    }

    async onMessage(data: ArrayBuffer): Promise<void> {
        const frame = GrpcFrame.fromBinary(new Uint8Array(data));
        const streamId = frame.streamId;
        const stream = this.activeStreams.get(streamId);

        if (stream) {
            switch (frame.payload.oneofKind) {
                case 'header': {
                    stream[0].debug && console.debug('GRPC-WS: received header for stream', streamId);

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
                    stream[0].debug && console.debug('GRPC-WS: received body for stream', streamId);

                    const body = frame.payload.body;
                    if (body === null) {
                        return;
                    }

                    stream[0].onChunk(body.data);
                    break;
                }

                case 'complete': {
                    stream[0].debug && console.debug('GRPC-WS: received complete for stream', streamId);

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
                        console.debug(
                            'GRPC-WS: received failure for stream',
                            streamId,
                            failure.errorStatus,
                            failure.errorMessage,
                        );

                    const metaData = headersToMetadata(failure.headers);
                    stream[0].onEnd(createRpcError(metaData, stream[0].methodDefinition));
                    break;
                }

                case 'cancel': {
                    stream[0].debug && console.debug('GRPC-WS: received cancel for stream', streamId);

                    stream[0].onEnd(errCancelled);
                    break;
                }

                default:
                    stream[0].debug &&
                        console.debug('GRPC-WS: received unknown message type for stream', streamId, frame.payload.oneofKind);
                    break;
            }
        } else {
            console.warn('GRPC-WS: stream does not exist', streamId);
        }
    }

    getStream(opts: TransportOptions): GrpcStream {
        let currentStreamId = this.streamId++;
        const self = this;

        async function sendToWebsocket(toSend: GrpcFrame): Promise<void> {
            if (!self.activeStreams.has(toSend.streamId)) {
                opts.debug && console.debug('GRPC-WS: stream does not exist', toSend.streamId);
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

        //question: can this structure be reused or is it one time use?
        const stream = {
            streamId: currentStreamId,

            start: (metadata: Metadata) => {
                opts.debug &&
                    console.debug(
                        'GRPC-WS: stream start',
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
                opts.debug && console.debug('GRPC-WS: stream send', currentStreamId);

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
                opts.debug && console.debug('GRPC-WS: stream complete', currentStreamId);

                const frame = newFrame();
                frame.payload = {
                    oneofKind: 'complete',
                    complete: Complete.create(),
                };

                sendToWebsocket(frame);
            },

            cancel: async () => {
                opts.debug && console.debug('GRPC-WS: stream cancel', currentStreamId);

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

import { RpcError, type MethodInfo } from '@protobuf-ts/runtime-rpc';
import { Body, Cancel, Complete, GrpcFrame, Header, HeaderValue } from '~~/gen/ts/resources/common/grpcws/grpcws';
import { headersToMetadata } from '../../bridge/utils';
import { Metadata } from '../../metadata';
import { type Transport, type TransportFactory, type TransportOptions } from '../transport';

type WebsocketAddress = string;

const activeWebsockets = new Map<WebsocketAddress, WebsocketChannel>();

export function WebsocketChannelTransport(): TransportFactory {
    function getChannel(wsHost: string): WebsocketChannel {
        let channel = activeWebsockets.get(wsHost);
        if (channel == null) {
            let newChannel = new WebsocketChannelImpl(wsHost);
            activeWebsockets.set(wsHost, newChannel);
            return newChannel;
        } else {
            return channel;
        }
    }

    return (opts: TransportOptions) => {
        let wsHost = constructWebSocketAddress(opts.url);

        return getChannel(wsHost).websocketChannelRequest(opts);
    };
}

export function closeWebsocketChannels(): void {
    activeWebsockets.forEach((aws) => aws.close());

    activeWebsockets.clear();
    console.info('GRPC-WS: closed websocket');
}

interface GrpcStream extends Transport {
    readonly streamId: number;
    flush(): void;
}

interface WebsocketChannel {
    websocketChannelRequest(options: TransportOptions): GrpcStream;
    close(): void;
}

function createRpcError(metaData: Metadata, methodDefinition: MethodInfo<object, object>): RpcError | undefined {
    if (!metaData.has('grpc-message') || !metaData.has('grpc-status')) {
        return;
    }

    const status = metaData.get('grpc-status').at(0);
    if (status === '0') {
        return;
    }
    const message = metaData.get('grpc-message').at(0) ?? '';

    const err = new RpcError(message, status ?? '0', metaData.headersMap);
    err.serviceName = methodDefinition.service.typeName;
    err.methodName = methodDefinition.name;
    return err;
}

class WebsocketChannelImpl implements WebsocketChannel {
    readonly wsUrl: string;
    readonly activeStreams = new Map<number, [TransportOptions, GrpcStream]>();
    ws: WebSocket | null = null;
    streamId = 0;
    closed = false;

    constructor(ws: string) {
        this.wsUrl = ws;
    }

    flush(event: Event) {
        this.activeStreams.forEach((opts, streamId, _map) => {
            opts[0].debug && console.debug('GRPC-WS: channel opened', streamId, event);
            opts[1].flush();
        });
        this.activeStreams;
    }

    recover(event: CloseEvent | Event) {
        if (this.closed) {
            return;
        }

        if (event instanceof CloseEvent) {
            this.activeStreams.forEach((opts, streamId, _map) => {
                opts[0].debug && console.debug('GRPC-WS: channel close event', streamId, event.reason, event.code);

                opts[0].onEnd(event.code === 1001 ? event : undefined);
            });
            this.activeStreams.clear();

            // Restart websocket when it isn't "going away" or the browser is navigating away
            if (event.code !== 1001) {
                this.ws = null;
                this.getWebsocket();
            }
        } else {
            this.activeStreams.forEach((opts, streamId, _map) => {
                opts[0].debug && console.debug('GRPC-WS: channel error', streamId, event);

                opts[0].onEnd(event);
            });
            this.activeStreams.clear();

            this.ws = null;
            this.getWebsocket();
        }
    }

    close() {
        this.closed = true;
        this.ws?.close();
    }

    onMessage(event: MessageEvent) {
        const frame = GrpcFrame.fromBinary(new Uint8Array(event.data));
        const streamId = frame.streamId;
        console.debug('GRPC-WS: message for stream received', streamId);

        const stream = this.activeStreams.get(streamId);
        if (stream) {
            switch (frame.payload.oneofKind) {
                case 'header':
                    stream[0].debug && console.debug('GRPC-WS: received header for stream', streamId);
                    const header = frame.payload.header;
                    if (header !== null) {
                        const metaData = headersToMetadata(header.headers);
                        stream[0].onHeaders(metaData, header.status);

                        const err = createRpcError(metaData, stream[0].methodDefinition);
                        if (err) {
                            stream[0].onEnd(err);
                        }
                    }
                    break;

                case 'body':
                    stream[0].debug && console.debug('GRPC-WS: received body for', streamId);
                    const body = frame.payload.body;
                    if (body !== null) {
                        stream[0].onChunk(body.data);
                    }
                    break;

                case 'complete':
                    stream[0].debug && console.debug('GRPC-WS: completing', streamId);
                    stream[0].onEnd();
                    break;

                case 'failure':
                    const failure = frame.payload.failure;
                    if (failure !== null) {
                        const message = failure.errorMessage;
                        stream[0].debug && console.debug('GRPC-WS: failed', streamId, message);
                        const metaData = headersToMetadata(failure.headers);
                        stream[0].onEnd(createRpcError(metaData, stream[0].methodDefinition));
                    } else {
                        stream[0].onEnd(new Error('GRPC-WS: unknown error'));
                    }
                    break;

                case 'cancel':
                    stream[0].onEnd(new Error('GRPC-WS: stream was canceled'));
                    break;

                default:
                    break;
            }
        } else {
            // TODO better logging?
            console.warn('GRPC-WS: stream does not exist', streamId);
        }
    }

    private getWebsocket(): WebSocket {
        if (this.ws === null) {
            this.ws = new WebSocket(this.wsUrl, ['grpc-websocket-channel']);
            this.ws.binaryType = 'arraybuffer';
            this.ws.onopen = (event) => this.flush(event);
            this.ws.onclose = (event: CloseEvent) => this.recover(event);
            this.ws.onerror = (event) => this.recover(event);
            this.ws.onmessage = (event) => this.onMessage(event);
        }
        return this.ws;
    }

    websocketChannelRequest(opts: TransportOptions): GrpcStream {
        let currentStreamId = this.streamId++;
        let sendQueue: Array<GrpcFrame> = [];
        const self = this;

        async function sendToWebsocket(toSend: GrpcFrame): Promise<void> {
            if (self.activeStreams.get(toSend.streamId) !== null) {
                const ws = self.getWebsocket();
                if (ws.readyState === ws.CONNECTING) {
                    opts.debug && console.debug('GRPC-WS: stream.webscocket queued', currentStreamId, toSend.payload.oneofKind);
                    sendQueue.push(toSend);
                } else {
                    sendQueue.forEach((toSend) => {
                        console.debug('GRPC-WS: sending from queue', currentStreamId, toSend.payload.oneofKind);
                        ws.send(GrpcFrame.toBinary(toSend));
                    });
                    sendQueue = [];
                    ws.send(GrpcFrame.toBinary(toSend));
                }
            } else {
                console.debug('GRPC-WS: stream does not exist', toSend.streamId);
            }
        }

        function newFrame(): GrpcFrame {
            const frame = GrpcFrame.create();
            frame.streamId = currentStreamId;
            return frame;
        }

        function writeUInt32BE(arr: Uint8Array, value: number, offset: number) {
            value = +value;
            offset = offset | 0;
            arr[offset] = value >>> 24;
            arr[offset + 1] = value >>> 16;
            arr[offset + 2] = value >>> 8;
            arr[offset + 3] = value & 0xff;
            return offset + 4;
        }

        //question: can this structure be reused or is it one time use?
        const stream = {
            streamId: currentStreamId,
            start: (metadata: Metadata) => {
                opts.debug &&
                    console.debug(
                        'GRPC-WS: stream.start',
                        currentStreamId,
                        `${opts.methodDefinition.service.typeName}/${opts.methodDefinition.name}`,
                    );
                self.activeStreams.set(currentStreamId, [opts, stream]);

                const header = Header.create();
                header.operation = `${opts.methodDefinition.service.typeName}/${opts.methodDefinition.name}`;
                // TODO add all meta data
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
                opts.debug && console.debug('GRPC-WS: stream.sendMessage', currentStreamId);
                const body = Body.create();
                const output = new Uint8Array(msgBytes.length + 5);
                output[0] = 0; // Compression none
                writeUInt32BE(output, msgBytes.length, 1);
                output.set(msgBytes, 5);
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
                opts.debug && console.debug('GRPC-WS: stream.finished', currentStreamId);
                const frame = newFrame();
                frame.payload = {
                    oneofKind: 'complete',
                    complete: Complete.create(),
                };

                sendToWebsocket(frame);
            },
            cancel: async () => {
                opts.debug && console.debug('GRPC-WS: stream.abort', currentStreamId);
                const frame = newFrame();
                frame.payload = {
                    oneofKind: 'cancel',
                    cancel: Cancel.create(),
                };

                await sendToWebsocket(frame);
            },
            flush: () => {
                opts.debug && console.debug('GRPC-WS: stream.flushed', currentStreamId);
                const ws = self.getWebsocket();
                if (ws.readyState === ws.OPEN) {
                    sendQueue.forEach((toSend) => {
                        console.debug('GRPC-WS: sending from queue', currentStreamId, toSend.payload.oneofKind);
                        ws.send(GrpcFrame.toBinary(toSend));
                    });
                    sendQueue = [];
                }
            },
        };
        return stream;
    }
}

function constructWebSocketAddress(url: string) {
    if (url.substr(0, 8) === 'https://') {
        return `wss://${url.substr(8)}`;
    } else if (url.substr(0, 7) === 'http://') {
        return `ws://${url.substr(7)}`;
    }
    throw new Error('Websocket transport constructed with non-https:// or http:// host.' + url);
}

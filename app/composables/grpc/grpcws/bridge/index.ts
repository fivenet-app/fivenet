import { GrpcStatusCode } from '@protobuf-ts/grpcweb-transport';
import {
    ClientStreamingCall,
    Deferred,
    DeferredState,
    DuplexStreamingCall,
    RpcError,
    RpcOutputStreamController,
    ServerStreamingCall,
    mergeRpcOptions,
    type MethodInfo,
    type RpcInputStream,
    type RpcMetadata,
    type RpcOptions,
    type RpcStatus,
    type RpcTransport,
    type UnaryCall,
} from '@protobuf-ts/runtime-rpc';
import type { UseWebSocketReturn } from '@vueuse/core';
import { Metadata } from '~/composables/grpc/grpcws/metadata';
import { GrpcFrame } from '~~/gen/ts/resources/common/grpcws/grpcws';
import type { GrpcWSOptions } from '../../grpcws/bridge/options';
import { errInternal, errTimeout, errUnavailable } from '../errors';
import type { Transport, TransportFactory } from '../transports/transport';
import { WebsocketChannelTransport } from '../transports/websocket/websocketChannel';
import { constructWebSocketAddress, createGrpcStatus, createGrpcTrailers } from './utils';

const logger = useLogger('📡 GRPC-WS');

const pingBinaryMsg = GrpcFrame.toBinary({
    streamId: 0,
    payload: {
        oneofKind: 'ping',
        ping: {
            pong: false,
        },
    },
});

const heartbeatMsg = pingBinaryMsg.buffer.slice(
    pingBinaryMsg.byteOffset,
    pingBinaryMsg.byteOffset + pingBinaryMsg.byteLength,
) as ArrayBuffer;

export const webSocket = useWebSocket(
    constructWebSocketAddress(
        `${window.location.protocol}//${window.location.hostname}:${!import.meta.dev ? window.location.port : 8080}/api/grpcws`,
    ),
    {
        immediate: false,
        autoReconnect: {
            delay: 350,
        },
        protocols: ['grpc-websocket-channel'],
        heartbeat: {
            message: heartbeatMsg,
            interval: 35_000,
            pongTimeout: 1_500,
        },

        onConnected(ws) {
            ws.binaryType = 'arraybuffer';
            logger.info('Websocket connected');
        },
        onDisconnected(_, event) {
            if (event.wasClean) {
                logger.info('Websocket disconnected cleanly, code:', event.code, 'reason:', event.reason);
                return;
            }

            logger.error('Websocket disconnected, code:', event.code, 'reason:', event.reason);
        },
    },
);

export class GrpcWSTransport implements RpcTransport {
    private readonly defaultOptions;
    webSocket: UseWebSocketReturn<ArrayBuffer>;
    private wsTs: TransportFactory;

    constructor(defaultOptions: GrpcWSOptions) {
        this.defaultOptions = defaultOptions;

        this.webSocket = webSocket;
        this.wsTs = WebsocketChannelTransport(logger, this.webSocket);
    }

    mergeOptions(options?: Partial<RpcOptions>): RpcOptions {
        return mergeRpcOptions(this.defaultOptions, options);
    }

    unary<I extends object, O extends object>(method: MethodInfo<I, O>, _input: I, _options: RpcOptions): UnaryCall<I, O> {
        const e = new RpcError('Unary request is not supported by grpc-ws', GrpcStatusCode[GrpcStatusCode.UNIMPLEMENTED]);
        e.methodName = method.name;
        e.serviceName = method.service.typeName;
        throw e;
    }

    serverStreaming<I extends object, O extends object>(
        method: MethodInfo<I, O>,
        input: I,
        options: RpcOptions,
    ): ServerStreamingCall<I, O> {
        if (this.webSocket.status.value !== 'OPEN') {
            logger.error("Websocket isn't connected, cannot create server streaming call");
            throw errUnavailable;
        }

        const opt = options as GrpcWSOptions,
            transport = this.wsTs({
                methodDefinition: method,
                debug: opt.debug,
                url: '',

                onChunk(chunkBytes) {
                    if (outStream.closed) {
                        transport.cancel();
                        return;
                    }

                    outStream.notifyMessage(method.O.fromBinary(chunkBytes, opt.binaryOptions));
                },
                onEnd(err) {
                    if (err !== undefined && !(err instanceof RpcError)) {
                        if (err.name === 'AbortError') {
                            err = errTimeout;
                        } else {
                            err = errInternal;
                        }
                    }

                    defHeader.rejectPending(err);
                    defStatus.rejectPending(err);
                    if (!outStream.closed) {
                        if (err) {
                            outStream.notifyError(err);
                        } else {
                            outStream.notifyComplete();
                        }
                    }
                    defTrailer.rejectPending(err);
                },
                onHeaders(headers: Metadata, _: number): void {
                    defHeader.resolvePending(headers.headersMap);

                    defStatus.resolvePending(createGrpcStatus(headers));
                    defTrailer.resolvePending(createGrpcTrailers(headers));
                },
            });

        const meta = opt.meta ?? {},
            defHeader = new Deferred<RpcMetadata>(),
            outStream = new RpcOutputStreamController<O>(),
            defStatus = new Deferred<RpcStatus>(),
            defTrailer = new Deferred<RpcMetadata>(),
            call = new ServerStreamingCall<I, O>(
                method,
                meta,
                input,
                defHeader.promise,
                outStream,
                defStatus.promise,
                defTrailer.promise,
            );

        let timeoutId: NodeJS.Timeout | undefined;

        if (opt.abort) {
            opt.abort.addEventListener('abort', (_) => {
                if (timeoutId) clearTimeout(timeoutId);
                transport.cancel();
            });
        }

        transport.start(new Metadata());
        transport.sendMessage(method.I.toBinary(input, opt.binaryOptions), true);

        return call;
    }

    clientStreaming<I extends object, O extends object>(
        method: MethodInfo<I, O>,
        options: RpcOptions,
    ): ClientStreamingCall<I, O> {
        if (this.webSocket.status.value !== 'OPEN') {
            logger.error("Websocket isn't connected, cannot create client streaming call");
            throw errUnavailable;
        }

        const opt = options as GrpcWSOptions,
            transport = this.wsTs({
                methodDefinition: method,
                debug: opt.debug,
                url: '',

                onChunk(chunkBytes) {
                    defMessage.resolve(method.O.fromBinary(chunkBytes, opt.binaryOptions));
                },
                onEnd(err) {
                    if (err !== undefined && !(err instanceof RpcError)) {
                        if (err.name === 'AbortError') {
                            err = errTimeout;
                        } else {
                            err = errInternal;
                        }
                    }

                    defHeader.rejectPending(err);
                    defMessage.rejectPending(err);
                    defStatus.rejectPending(err);
                    defTrailer.rejectPending(err);

                    if (defMessage.state === DeferredState.RESOLVED) {
                        return;
                    }
                    defMessage.rejectPending(err);
                },
                onHeaders(headers: Metadata, _: number): void {
                    defHeader.resolvePending(headers.headersMap);

                    defStatus.resolvePending(createGrpcStatus(headers));
                    defTrailer.resolvePending(createGrpcTrailers(headers));
                },
            });

        const meta = opt.meta ?? {},
            defHeader = new Deferred<RpcMetadata>(),
            defMessage = new Deferred<O>(),
            defStatus = new Deferred<RpcStatus>(),
            defTrailer = new Deferred<RpcMetadata>(),
            inStream = new GrpcInputStreamWrapper(transport, (m) => method.I.toBinary(m as I, opt.binaryOptions)),
            call = new ClientStreamingCall<I, O>(
                method,
                meta,
                inStream,
                defHeader.promise,
                defMessage.promise,
                defStatus.promise,
                defTrailer.promise,
            );

        let timeoutId: NodeJS.Timeout | undefined;

        if (opt.abort) {
            opt.abort.addEventListener('abort', (_) => {
                if (timeoutId) clearTimeout(timeoutId);
                transport.cancel();
            });
        }

        transport.start(new Metadata());

        return call;
    }

    duplex<I extends object, O extends object>(method: MethodInfo<I, O>, options: RpcOptions): DuplexStreamingCall<I, O> {
        if (this.webSocket.status.value !== 'OPEN') {
            logger.error("Websocket isn't connected, cannot create duplex streaming call");
            throw errUnavailable;
        }

        const opt = options as GrpcWSOptions,
            transport = this.wsTs({
                methodDefinition: method,
                debug: opt.debug,
                url: '',

                onChunk(chunkBytes) {
                    outStream.notifyMessage(method.O.fromBinary(chunkBytes, opt.binaryOptions));
                },
                onEnd(err) {
                    if (err !== undefined && !(err instanceof RpcError)) {
                        if (err.name === 'AbortError') {
                            err = errTimeout;
                        } else {
                            err = errInternal;
                        }
                    }

                    defHeader.rejectPending(err);
                    defStatus.rejectPending(err);
                    if (!outStream.closed) {
                        if (err) {
                            outStream.notifyError(err);
                        } else {
                            outStream.notifyComplete();
                        }
                    }
                    defTrailer.rejectPending(err);
                },
                onHeaders(headers: Metadata, _: number): void {
                    defHeader.resolvePending(headers.headersMap);

                    defStatus.resolvePending(createGrpcStatus(headers));
                    defTrailer.resolvePending(createGrpcTrailers(headers));
                },
            });

        const meta = opt.meta ?? {},
            defHeader = new Deferred<RpcMetadata>(),
            outStream = new RpcOutputStreamController<O>(),
            defStatus = new Deferred<RpcStatus>(),
            defTrailer = new Deferred<RpcMetadata>(),
            inStream = new GrpcInputStreamWrapper(transport, (m) => method.I.toBinary(m as I, opt.binaryOptions)),
            call = new DuplexStreamingCall<I, O>(
                method,
                meta,
                inStream,
                defHeader.promise,
                outStream,
                defStatus.promise,
                defTrailer.promise,
            );

        let timeoutId: NodeJS.Timeout | undefined;

        if (opt.abort) {
            opt.abort.addEventListener('abort', (_) => {
                if (timeoutId) clearTimeout(timeoutId);
                transport.cancel();
            });
        }

        transport.start(new Metadata());

        return call;
    }

    close(): void {
        if (this.webSocket.status.value === 'CLOSED') return;

        logger.info('Closing Websocket');
        this.webSocket.close();
    }
}

class GrpcInputStreamWrapper<I> implements RpcInputStream<I> {
    protected toBinary: (message: I) => Uint8Array;

    constructor(
        private readonly transport: Transport,
        toBinary: (message: I) => Uint8Array,
    ) {
        this.toBinary = toBinary;
    }

    send(message: I): Promise<void> {
        return new Promise<void>((resolve, _) => {
            this.transport.sendMessage(this.toBinary(message));
            resolve();
        });
    }

    complete(): Promise<void> {
        this.transport.finishSend();
        return Promise.resolve(undefined);
    }
}

import {
    ClientStreamingCall,
    Deferred,
    DuplexStreamingCall,
    RpcOutputStreamController,
    ServerStreamingCall,
    UnaryCall,
    mergeRpcOptions,
    type MethodInfo,
    type RpcInputStream,
    type RpcMetadata,
    type RpcOptions,
    type RpcStatus,
    type RpcTransport,
} from '@protobuf-ts/runtime-rpc';
import type { UseWebSocketReturn } from '@vueuse/core';
import { Metadata } from '~/composables/grpcws/metadata';
import type { GrpcWSOptions } from '../../grpcws/bridge/options';
import type { Transport, TransportFactory } from '../transports/transport';
import { WebsocketChannelTransport } from '../transports/websocket/websocketChannel';
import { createGrpcStatus, createGrpcTrailers } from './utils';

export class GrpcWSTransport implements RpcTransport {
    private readonly defaultOptions;
    private logger: ILogger;
    webSocket: UseWebSocketReturn<any>;
    wsInitiated: Ref<boolean>;
    private wsTs: TransportFactory;

    constructor(defaultOptions: GrpcWSOptions) {
        this.defaultOptions = defaultOptions;

        this.logger = useLogger('📡 GRPC-WS');

        this.wsInitiated = ref(false);
        const self = this;
        const webSocket = useWebSocket(defaultOptions.wsUrl, {
            immediate: false,
            autoReconnect: {
                delay: 1150,
            },
            protocols: ['grpc-websocket-channel'],
            onConnected(ws) {
                ws.binaryType = 'arraybuffer';
                self.wsInitiated.value = true;
                self.logger.info('Websocket connected');
            },
            onDisconnected(_, event) {
                if (event.wasClean) {
                    return;
                }

                self.logger.error('Websocket disconnected, code:', event.code, 'reason:', event.reason);
            },
        });
        this.webSocket = webSocket;
        this.wsTs = WebsocketChannelTransport(this.logger, this.webSocket);
    }

    mergeOptions(options?: Partial<RpcOptions>): RpcOptions {
        return mergeRpcOptions(this.defaultOptions, options);
    }

    unary<I extends object, O extends object>(method: MethodInfo<I, O>, input: I, options: RpcOptions): UnaryCall<I, O> {
        const opt = options as GrpcWSOptions,
            transport = this.wsTs({
                methodDefinition: method,
                debug: opt.debug,
                url: '',

                onChunk(chunkBytes) {
                    defHeader.resolvePending({});
                    defTrailer.resolvePending({});
                    defStatus.resolvePending(createGrpcStatus(new Metadata()));
                    defMessage.resolve(method.O.fromBinary(chunkBytes, opt.binaryOptions));
                },
                onEnd(err) {
                    if (err) {
                        defHeader.rejectPending(err);
                        defMessage.rejectPending(err);
                        defStatus.rejectPending(err);
                        defTrailer.rejectPending(err);
                        defMessage.rejectPending(err);
                        return;
                    }
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
            call = new UnaryCall<I, O>(
                method,
                meta,
                input,
                defHeader.promise,
                defMessage.promise,
                defStatus.promise,
                defTrailer.promise,
            );

        const abort = opt.abort || (opt.timeout ? AbortSignal.timeout(opt.timeout) : undefined);
        if (abort) {
            abort.addEventListener('abort', (_) => {
                transport.cancel();
            });
        }

        transport.start(new Metadata());
        transport.sendMessage(method.I.toBinary(input, opt.binaryOptions), true);

        return call;
    }

    serverStreaming<I extends object, O extends object>(
        method: MethodInfo<I, O>,
        input: I,
        options: RpcOptions,
    ): ServerStreamingCall<I, O> {
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
                    if (err && err instanceof Error) {
                        defHeader.rejectPending(err);
                        outStream.notifyError(err);
                        defStatus.rejectPending(err);
                        defTrailer.rejectPending(err);
                        return;
                    }

                    if (!outStream.closed) {
                        outStream.notifyComplete();
                    }
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

        if (opt.abort) {
            opt.abort.addEventListener('abort', (_) => {
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
        const opt = options as GrpcWSOptions,
            transport = this.wsTs({
                methodDefinition: method,
                debug: opt.debug,
                url: '',

                onChunk(chunkBytes) {
                    defMessage.resolve(method.O.fromBinary(chunkBytes, opt.binaryOptions));
                },
                onEnd(err) {
                    if (err) {
                        defHeader.rejectPending(err);
                        defMessage.rejectPending(err);
                        defStatus.rejectPending(err);
                        defTrailer.rejectPending(err);
                        return;
                    }

                    defMessage.resolve(method.O.create());
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

        if (opt.abort) {
            opt.abort.addEventListener('abort', (_) => {
                transport.cancel();
            });
        }

        transport.start(new Metadata());

        return call;
    }

    duplex<I extends object, O extends object>(method: MethodInfo<I, O>, options: RpcOptions): DuplexStreamingCall<I, O> {
        const opt = options as GrpcWSOptions,
            transport = this.wsTs({
                methodDefinition: method,
                debug: opt.debug,
                url: '',

                onChunk(chunkBytes) {
                    outStream.notifyMessage(method.O.fromBinary(chunkBytes, opt.binaryOptions));
                },
                onEnd(err) {
                    if (err) {
                        defHeader.rejectPending(err);
                        outStream.notifyError(err);
                        defStatus.rejectPending(err);
                        defTrailer.rejectPending(err);
                        return;
                    }

                    if (!outStream.closed) {
                        outStream.notifyComplete();
                    }
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

        if (opt.abort) {
            opt.abort.addEventListener('abort', (_) => {
                transport.cancel();
            });
        }

        transport.start(new Metadata());

        return call;
    }

    close(): void {
        if (this.webSocket.status.value === 'CLOSED') {
            return;
        }

        this.webSocket.close();
        this.logger.info('Closed Websocket');
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

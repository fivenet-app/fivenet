import type {
    ClientStreamingCall,
    DuplexStreamingCall,
    MethodInfo,
    RpcOptions,
    RpcTransport,
    ServerStreamingCall,
    UnaryCall,
} from '@protobuf-ts/runtime-rpc';

export class GrpcCombinedTransport {
    private unaryClient: RpcTransport;
    private streamClient: RpcTransport;

    constructor(unaryClient: RpcTransport, streamClient: RpcTransport) {
        this.unaryClient = unaryClient;
        this.streamClient = streamClient;
    }

    mergeOptions(options?: Partial<RpcOptions>): RpcOptions {
        return this.streamClient.mergeOptions(options);
    }

    unary<I extends object, O extends object>(method: MethodInfo<I, O>, input: I, options: RpcOptions): UnaryCall<I, O> {
        return this.unaryClient.unary<I, O>(method, input, options);
    }

    serverStreaming<I extends object, O extends object>(
        method: MethodInfo<I, O>,
        input: I,
        options: RpcOptions,
    ): ServerStreamingCall<I, O> {
        return this.streamClient.serverStreaming<I, O>(method, input, options);
    }

    clientStreaming<I extends object, O extends object>(
        method: MethodInfo<I, O>,
        options: RpcOptions,
    ): ClientStreamingCall<I, O> {
        return this.streamClient.clientStreaming<I, O>(method, options);
    }

    duplex<I extends object, O extends object>(method: MethodInfo<I, O>, options: RpcOptions): DuplexStreamingCall<I, O> {
        return this.streamClient.duplex<I, O>(method, options);
    }
}

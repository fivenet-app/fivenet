import { GrpcWebFetchTransport, type GrpcWebOptions } from '@protobuf-ts/grpcweb-transport';
import type {
    ClientStreamingCall,
    DuplexStreamingCall,
    MethodInfo,
    RpcOptions,
    RpcTransport,
    ServerStreamingCall,
    UnaryCall,
} from '@protobuf-ts/runtime-rpc';
import { useGRPCWebsocketTransport } from './grpcws';

// Lazy singleton instance
let _transport: GrpcCombinedTransport | null = null;

export function useGRPCTransport() {
    if (!_transport) {
        const options: GrpcWebOptions = {
            baseUrl: '/api/grpc',
            format: 'text',
            fetchInit: { credentials: 'same-origin' },
        };
        const grpcWebTransport = new GrpcWebFetchTransport(options);
        const grpcWebsocketTransport = useGRPCWebsocketTransport();
        _transport = new GrpcCombinedTransport(grpcWebTransport, grpcWebsocketTransport);
    }
    return _transport;
}

function authInterceptor(options: RpcOptions): RpcOptions {
    // Interceptros don't seem to work 100% of the time.. probably because of the "Frankenstein" transport setup.
    if (!options) options = {};
    if (!options.meta) options.meta = {};

    const userToken = sessionStorage.getItem(authUserTokenKey);
    if (userToken) {
        options.meta['Authorization'] = `Bearer ${userToken}`;
    }

    return options;
}

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
        // Interceptors don't seem to work 100% of the time (at least for unary calls)..
        // probably because of the "Frankenstein" transport setup.
        options = authInterceptor(options);
        return this.unaryClient.unary<I, O>(method, input, this.unaryClient.mergeOptions(options));
    }

    serverStreaming<I extends object, O extends object>(
        method: MethodInfo<I, O>,
        input: I,
        options: RpcOptions,
    ): ServerStreamingCall<I, O> {
        options = authInterceptor(options);
        return this.streamClient.serverStreaming<I, O>(method, input, options);
    }

    clientStreaming<I extends object, O extends object>(
        method: MethodInfo<I, O>,
        options: RpcOptions,
    ): ClientStreamingCall<I, O> {
        options = authInterceptor(options);
        return this.streamClient.clientStreaming<I, O>(method, options);
    }

    duplex<I extends object, O extends object>(method: MethodInfo<I, O>, options: RpcOptions): DuplexStreamingCall<I, O> {
        options = authInterceptor(options);
        return this.streamClient.duplex<I, O>(method, options);
    }
}

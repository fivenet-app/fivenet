import { GrpcWebFetchTransport, type GrpcWebOptions } from '@protobuf-ts/grpcweb-transport';
import type { RpcInterceptor } from '@protobuf-ts/runtime-rpc';
import { GrpcCombinedTransport } from './grpc/combined';
import { useGRPCWebsocketTransport } from './grpc/grpcws';

// Lazy singleton instance
let _transport: GrpcCombinedTransport | null = null;

export function useGRPCTransport() {
    if (!_transport) {
        const options: GrpcWebOptions = {
            baseUrl: '/api/grpc',
            format: 'text',
            fetchInit: { credentials: 'same-origin' },
            interceptors: [authInterceptor],
        };
        const grpcWebTransport = new GrpcWebFetchTransport(options);
        const grpcWebsocketTransport = useGRPCWebsocketTransport();
        _transport = new GrpcCombinedTransport(grpcWebTransport, grpcWebsocketTransport);
    }
    return _transport;
}

export const authInterceptor = {
    interceptUnary(next, method, input, options) {
        if (!options.meta) {
            options.meta = {};
        }

        const userToken = sessionStorage.getItem('fivenet:user_token_v1');
        if (userToken) {
            options.meta['Authorization'] = `Bearer ${userToken}`;
        }

        return next(method, input, options);
    },
} as RpcInterceptor;

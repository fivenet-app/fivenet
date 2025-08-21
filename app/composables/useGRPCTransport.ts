import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { GrpcCombinedTransport } from './grpc/combined';
import { useGRPCWebsocketTransport } from './grpc/grpcws';

// Lazy singleton instance
let _transport: GrpcCombinedTransport | null = null;

export function useGRPCTransport() {
    if (!_transport) {
        const grpcWebTransport = new GrpcWebFetchTransport({
            baseUrl: '/api/grpc',
            format: 'text',
            fetchInit: { credentials: 'same-origin' },
        });
        const grpcWebsocketTransport = useGRPCWebsocketTransport();
        _transport = new GrpcCombinedTransport(grpcWebTransport, grpcWebsocketTransport);
    }
    return _transport;
}

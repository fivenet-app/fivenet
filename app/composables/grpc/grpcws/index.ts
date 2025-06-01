import { GrpcWSTransport } from './bridge';

function _useGRPCWebsocketTransport(): GrpcWSTransport {
    const grpcWebsocketTransport = new GrpcWSTransport({
        debug: import.meta.dev,
        timeout: 8500,

        // GRPC web transport options
        baseUrl: '/api/grpc',
        format: 'text',
        fetchInit: {
            credentials: 'same-origin',
        },
    });

    return grpcWebsocketTransport;
}

export const useGRPCWebsocketTransport = createSharedComposable(_useGRPCWebsocketTransport);

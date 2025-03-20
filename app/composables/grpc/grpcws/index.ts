import { GrpcWSTransport } from './bridge';
import { constructWebSocketAddress } from './bridge/utils';

function _useGRPCWebsocketTransport(): GrpcWSTransport {
    const grpcWebsocketTransport = new GrpcWSTransport({
        wsUrl: constructWebSocketAddress(
            `${window.location.protocol}//${window.location.hostname}:${!import.meta.dev ? window.location.port : 8080}/api/grpcws`,
        ),
        debug: import.meta.dev,
        timeout: 8500,
        reconnect: true,

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

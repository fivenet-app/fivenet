import { GrpcWSTransport } from './bridge';
import { constructWebSocketAddress } from './bridge/utils';

function _useGRPCWebsocketTransport(): GrpcWSTransport {
    const grpcWebsocketTransport = new GrpcWSTransport({
        wsUrl: constructWebSocketAddress(
            `${window.location.protocol}//${window.location.hostname}:${!import.meta.dev ? window.location.port : 8080}/api/grpc`,
        ),
        debug: import.meta.dev,
        timeout: 8500,
        reconnect: true,
    });

    return grpcWebsocketTransport;
}

export const useGRPCWebsocketTransport = createSharedComposable(_useGRPCWebsocketTransport);

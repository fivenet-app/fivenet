import { GrpcWSTransport } from './bridge';
import { constructWebSocketAddress } from './bridge/utils';

let grpcWebsocketTransport: GrpcWSTransport | undefined = undefined;

export function useGRPCWebsocketTransport(): GrpcWSTransport {
    if (grpcWebsocketTransport === undefined) {
        grpcWebsocketTransport = new GrpcWSTransport({
            wsUrl: constructWebSocketAddress(
                `${window.location.protocol}//${window.location.hostname}:${!import.meta.dev ? window.location.port : 8080}/api/grpc`,
            ),
            debug: import.meta.dev,
            timeout: 8500,
        });
    }

    return grpcWebsocketTransport;
}

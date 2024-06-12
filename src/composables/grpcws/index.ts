import { GrpcWSTransport } from './bridge';
import { constructWebSocketAddress } from './bridge/utils';

let grpcWebsocketTransport: GrpcWSTransport | undefined = undefined;

export function useGRPCWebsocketTransport(): GrpcWSTransport {
    if (grpcWebsocketTransport === undefined) {
        grpcWebsocketTransport = new GrpcWSTransport({
            url: constructWebSocketAddress(
                `${window.location.protocol}//${window.location.hostname}:${!import.meta.dev ? window.location.port : 8080}/api/grpc`,
            ),
        });
    }

    return grpcWebsocketTransport;
}

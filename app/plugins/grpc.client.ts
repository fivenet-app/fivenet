import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { GRPCClients } from '~~/gen/ts/clients';
import { GrpcCombinedTransport } from '../composables/grpc/combined';
import { useGRPCWebsocketTransport } from '../composables/grpc/grpcws';

const grpcWebTransport = new GrpcWebFetchTransport({
    baseUrl: '/api/grpc',
    format: 'text',
    fetchInit: {
        credentials: 'same-origin',
    },
});

const grpcWebsocketTransport = useGRPCWebsocketTransport();

const grpcCombinedTransport = new GrpcCombinedTransport(grpcWebTransport, grpcWebsocketTransport);

export default defineNuxtPlugin({
    name: 'grpc',
    parallel: true,

    async setup(_) {
        return {
            provide: {
                grpc: new GRPCClients(grpcCombinedTransport),
            },
        };
    },
});

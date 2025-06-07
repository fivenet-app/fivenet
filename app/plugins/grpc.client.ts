import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { GrpcCombinedTransport } from '~/composables/grpc/combined';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { GRPCClients } from '~~/gen/ts/clients';

export default defineNuxtPlugin({
    name: 'grpc',
    parallel: true,

    async setup() {
        const grpcWebTransport = new GrpcWebFetchTransport({
            baseUrl: '/api/grpc',
            format: 'text',
            fetchInit: {
                credentials: 'same-origin',
            },
        });

        const grpcWebsocketTransport = useGRPCWebsocketTransport();

        const transport = new GrpcCombinedTransport(grpcWebTransport, grpcWebsocketTransport);

        return {
            provide: {
                grpc: new GRPCClients(transport),
            },
        };
    },
});

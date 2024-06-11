import { type RpcMetadata, type RpcStatus } from '@protobuf-ts/runtime-rpc';
import { Metadata } from '~/composables/grpcws/metadata';

export function createGrpcStatus(metaData: Metadata): RpcStatus {
    return {
        code: metaData?.get('grpc-status')?.at(0) ?? '0',
        detail: metaData?.get('grpc-message')?.at(0) ?? '',
    };
}

export function createGrpcTrailers(metaData: Metadata): RpcMetadata {
    const trailers = new Metadata();
    metaData.forEach((k, v) => {
        if (!k.startsWith('trailer+')) {
            return;
        }

        trailers.append(k, v);
    });

    return trailers.headersMap;
}

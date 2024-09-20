import type { RpcMetadata, RpcStatus } from '@protobuf-ts/runtime-rpc';
import { Metadata } from '~/composables/grpcws/metadata';
import type { HeaderValue } from '~~/gen/ts/resources/common/grpcws/grpcws';

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

export function headersToMetadata(headers: { [key: string]: HeaderValue }): Metadata {
    const metaData = new Metadata();
    Object.keys(headers).forEach((k) => metaData.append(k.replaceAll(':', '+'), headers[k]!.value));
    return metaData;
}

export function constructWebSocketAddress(url: string) {
    if (url.substr(0, 8) === 'https://') {
        return `wss://${url.substr(8)}`;
    } else if (url.substr(0, 7) === 'http://') {
        return `ws://${url.substr(7)}`;
    }
    throw new Error('Websocket transport constructed with non-https:// or http:// host.' + url);
}

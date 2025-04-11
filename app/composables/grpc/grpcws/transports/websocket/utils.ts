import { GrpcStatusCode } from '@protobuf-ts/grpcweb-transport';
import { RpcError, type MethodInfo } from '@protobuf-ts/runtime-rpc';
import type { Metadata } from '../../metadata';

export function createRpcError(metaData: Metadata, methodDefinition: MethodInfo<object, object>): RpcError | undefined {
    if (!metaData.has('grpc-message') || !metaData.has('grpc-status')) {
        return;
    }

    const status = metaData.get('grpc-status').at(0);
    if (status === '0') {
        return;
    }
    const message = metaData.get('grpc-message').at(0) ?? '';

    const err = new RpcError(message, status ?? '0', metaData.headersMap);
    err.serviceName = methodDefinition.service.typeName;
    err.methodName = methodDefinition.name;
    err.code = (GrpcStatusCode[parseInt(err.code)] ?? GrpcStatusCode[GrpcStatusCode.INTERNAL])?.toUpperCase();

    return err;
}

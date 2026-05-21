import { GrpcStatusCode } from '@protobuf-ts/grpcweb-transport';
import { RpcError, type MethodInfo } from '@protobuf-ts/runtime-rpc';
import type { Metadata } from '../../metadata';

export function createRpcError(metaData: Metadata, methodDefinition: MethodInfo<object, object>): RpcError | undefined {
    if (!metaData.has('grpc-message') || !metaData.has('grpc-status')) return;

    const status = metaData.get('grpc-status').at(0);
    if (status === '0') return;
    const message = metaData.get('grpc-message').at(0) ?? '';

    const err = new RpcError(message, status ?? '0', metaData.headersMap);
    err.serviceName = methodDefinition.service.typeName;
    err.methodName = methodDefinition.name;
    err.code = (GrpcStatusCode[parseInt(err.code)] ?? GrpcStatusCode[GrpcStatusCode.INTERNAL])?.toUpperCase();

    return err;
}

const retryableCodes = ['UNAVAILABLE', 'DEADLINE_EXCEEDED', 'INTERNAL', 'RESOURCE_EXHAUSTED'];

export function isRetryableError(err: unknown): boolean {
    if (err instanceof RpcError) {
        // Check if the error is retryable
        return retryableCodes.includes(err.code);
    }
    // If not an RpcError, we assume it's not retryable
    return false;
}

export function writeUInt32BE(arr: Uint8Array, value: number, offset: number): number {
    value = +value;
    offset = offset | 0;
    arr[offset] = value >>> 24;
    arr[offset + 1] = value >>> 16;
    arr[offset + 2] = value >>> 8;
    arr[offset + 3] = value & 0xff;
    return offset + 4;
}

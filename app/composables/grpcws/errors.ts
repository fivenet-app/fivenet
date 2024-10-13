import { RpcError } from '@protobuf-ts/runtime-rpc';

export const errTimeout = new RpcError('GRPC-WS: timeout reached', 'DEADLINE_EXCEEDED');
export const errCancelled = new RpcError('GRPC-WS: request cancelled', 'CANCELLED');
export const errInternal = new RpcError(
    'notifications.grpc_errors.unavailable.title;notifications.grpc_errors.unavailable.content',
    'INTERNAL',
);
export const errUnavailable = new RpcError('GRPC-WS: unavailable', 'UNAVAILABLE');

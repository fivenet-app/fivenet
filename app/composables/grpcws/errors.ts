import { RpcError } from '@protobuf-ts/runtime-rpc';

export const errTimeout = new RpcError(
    'notifications.grpc_errors.deadline_exceeded.title;notifications.grpc_errors.deadline_exceeded.content',
    'DEADLINE_EXCEEDED',
);
export const errCancelled = new RpcError(
    'notifications.grpc_errors.cancelled.title;notifications.grpc_errors.cancelled.content',
    'CANCELLED',
);
export const errInternal = new RpcError(
    'notifications.grpc_errors.internal.title;notifications.grpc_errors.internal.content',
    'INTERNAL',
);
export const errUnavailable = new RpcError(
    'notifications.grpc_errors.unavailable.title;notifications.grpc_errors.unavailable.content',
    'UNAVAILABLE',
);

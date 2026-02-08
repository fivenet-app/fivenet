import { RpcError } from '@protobuf-ts/runtime-rpc';

export const errTimeout = new RpcError(
    '{"title":{"key":"notifications.grpc_errors.deadline_exceeded.title"},"content":{"key":notifications.grpc_errors.deadline_exceeded.content"}}',
    'DEADLINE_EXCEEDED',
);
export const errCancelled = new RpcError(
    '{"title":{"key":"notifications.grpc_errors.cancelled.title"},"content":{"key":"notifications.grpc_errors.cancelled.content"}}',
    'CANCELLED',
);
export const errInternal = new RpcError(
    '{"title":{"key":"notifications.grpc_errors.internal.title"},"content":{"key":"notifications.grpc_errors.internal.content"}}',
    'INTERNAL',
);
export const errUnavailable = new RpcError(
    '{"title":{"key":"notifications.grpc_errors.unavailable.title"},"content":{"key":"notifications.grpc_errors.unavailable.content"}}',
    'UNAVAILABLE',
);

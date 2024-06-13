import { RpcError } from '@protobuf-ts/runtime-rpc';

export const errCancelled = new RpcError('GRPC-WS: timeout reached', 'CANCELLED');
export const errInternal = new RpcError('notifications.grpc_errors.unavailable.title;notifications.grpc_errors.unavailable.content', 'INTERNAL');
export const errUnavailable = new RpcError('GRPC-WS: unavailable', 'UNAVAILABLE');

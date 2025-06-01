import type { GrpcWebOptions } from '@protobuf-ts/grpcweb-transport';
import type { RpcOptions } from '@protobuf-ts/runtime-rpc';

export interface GrpcWSOptions extends RpcOptions, GrpcWebOptions {
    debug?: boolean;
}

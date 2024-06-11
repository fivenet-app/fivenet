import { type RpcOptions } from '@protobuf-ts/runtime-rpc';

export interface GrpcWSOptions extends RpcOptions {
    debug?: boolean;
    url: string;
}

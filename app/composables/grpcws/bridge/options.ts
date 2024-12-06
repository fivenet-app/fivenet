import type { RpcOptions } from '@protobuf-ts/runtime-rpc';

export interface GrpcWSOptions extends RpcOptions {
    debug?: boolean;
    wsUrl: string;
    timeout?: number;
    reconnect?: boolean;
}

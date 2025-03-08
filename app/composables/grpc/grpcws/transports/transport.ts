import type { MethodInfo } from '@protobuf-ts/runtime-rpc';
import type { Metadata } from '../metadata';

export interface Transport {
    sendMessage(msgBytes: Uint8Array, complete?: boolean): Promise<void>;
    finishSend(): Promise<void>;
    cancel(err?: Error): Promise<void>;
    start(metadata: Metadata): void;
}

export interface TransportOptions {
    methodDefinition: MethodInfo<object, object>;
    debug?: boolean;
    url: string;
    onHeaders: (headers: Metadata, status: number) => void;
    onChunk: (chunkBytes: Uint8Array, flush?: boolean) => void;
    onEnd: (err?: Error) => void;
}

export interface TransportFactory {
    (options: TransportOptions): Transport;
}

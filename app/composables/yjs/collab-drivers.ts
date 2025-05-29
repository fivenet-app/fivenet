import type { DuplexStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import type { ClientPacket, ServerPacket } from '~~/gen/ts/resources/collab/collab';

export type StreamConnectFn = (options?: RpcOptions) => DuplexStreamingCall<ClientPacket, ServerPacket>;

export type CollabCategory = 'documents' | 'wiki';

export const collabDrivers: Record<CollabCategory, StreamConnectFn> = {
    documents: (opts) => {
        const { $grpc } = useNuxtApp();
        return $grpc.documents.collab.joinRoom(opts);
    },
    wiki: (opts) => {
        const { $grpc } = useNuxtApp();
        return $grpc.wiki.collab.joinRoom(opts);
    },
};

import type { DuplexStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { getDocumentsCollabClient, getWikiCollabClient } from '~~/gen/ts/clients';
import type { ClientPacket, ServerPacket } from '~~/gen/ts/resources/collab/collab';

export type StreamConnectFn = (options?: RpcOptions) => DuplexStreamingCall<ClientPacket, ServerPacket>;

export type CollabCategory = 'documents' | 'wiki';

export const collabDrivers: Record<CollabCategory, () => Promise<StreamConnectFn>> = {
    documents: async () => {
        const client = await getDocumentsCollabClient();
        return (opts) => client.joinRoom(opts);
    },

    wiki: async () => {
        const client = await getWikiCollabClient();
        return (opts) => client.joinRoom(opts);
    },
};

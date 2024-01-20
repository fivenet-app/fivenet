import { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore, type StoreDefinition } from 'pinia';
import { ListColleaguesRequest, ListColleaguesResponse } from '~~/gen/ts/services/jobs/jobs';

export interface CompletorState {}

export const useJobsStore = defineStore('jobs', {
    state: () => ({}) as CompletorState,
    persist: false,
    actions: {
        // Colleagues
        async listColleagues(req: ListColleaguesRequest): Promise<ListColleaguesResponse> {
            const { $grpc } = useNuxtApp();

            try {
                const call = $grpc.getJobsClient().listColleagues(req);
                const { response } = await call;

                return response;
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useJobsStore as unknown as StoreDefinition, import.meta.hot));
}

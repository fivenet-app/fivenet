import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { Job } from '~~/gen/ts/resources/users/jobs';
import { CompleteCitizensRequest, CompleteJobsRequest } from '~~/gen/ts/services/completor/completor';
import { UserShort } from '../../gen/ts/resources/users/users';

export interface CompletorState {
    jobs: Job[];
}

export interface ClipboardState extends CompletorState {
    activeStack: CompletorState;
}

export const useCompletorStore = defineStore('completor', {
    state: () =>
        ({
            jobs: [] as Job[],
        }) as CompletorState,
    persist: false,
    actions: {
        // Jobs
        async getJobByName(name: string): Promise<Job | undefined> {
            return this.listJobs().then((jobs) => jobs.find((j) => j.name === name));
        },
        async listJobs(): Promise<Job[]> {
            if (this.jobs.length > 0) return this.jobs;

            this.jobs = await this.completeJobs({});
            return this.jobs;
        },

        // Citizens
        async findCitizen(userId: number): Promise<UserShort | undefined> {
            return this.completeCitizens({
                search: '',
                userId: userId,
            }).then((users) => (users.length === 0 ? undefined : users[0]));
        },

        // GRPC calls
        async completeJobs(req: CompleteJobsRequest): Promise<Job[]> {
            return new Promise(async (res, rej) => {
                const { $grpc } = useNuxtApp();
                try {
                    const call = $grpc.getCompletorClient().completeJobs(req);
                    const { response } = await call;

                    return res(response.jobs);
                } catch (e) {
                    $grpc.handleError(e as RpcError);
                    return rej(e as RpcError);
                }
            });
        },
        async completeCitizens(req: CompleteCitizensRequest): Promise<UserShort[]> {
            return new Promise(async (res, rej) => {
                const { $grpc } = useNuxtApp();
                try {
                    const call = $grpc.getCompletorClient().completeCitizens(req);
                    const { response } = await call;

                    return res(response.users);
                } catch (e) {
                    $grpc.handleError(e as RpcError);
                    return rej(e as RpcError);
                }
            });
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCompletorStore as unknown as StoreDefinition, import.meta.hot));
}

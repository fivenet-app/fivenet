import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { Category } from '~~/gen/ts/resources/documents/category';
import { LawBook } from '~~/gen/ts/resources/laws/laws';
import { Job } from '~~/gen/ts/resources/users/jobs';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { CompleteCitizensRequest, CompleteJobsRequest } from '~~/gen/ts/services/completor/completor';

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

        // Citizens
        async findCitizen(userId: number): Promise<UserShort | undefined> {
            return this.completeCitizens({
                search: '',
                userId: userId,
            }).then((users) => (users.length === 0 ? undefined : users[0]));
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

        // Document Categories
        async completeDocumentCategories(search: string): Promise<Category[]> {
            return new Promise(async (res, rej) => {
                if (!can('CompletorService.CompleteDocumentCategories')) {
                    return res([]);
                }

                const { $grpc } = useNuxtApp();
                try {
                    const call = $grpc.getCompletorClient().completeDocumentCategories({
                        search: search,
                    });
                    const { response } = await call;

                    return res(response.categories);
                } catch (e) {
                    $grpc.handleError(e as RpcError);
                    return rej(e as RpcError);
                }
            });
        },

        // Laws
        async listLawBooks(): Promise<LawBook[]> {
            return new Promise(async (res, rej) => {
                const { $grpc } = useNuxtApp();
                try {
                    const call = $grpc.getCompletorClient().listLawBooks({});
                    const { response } = await call;

                    return res(response.books);
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

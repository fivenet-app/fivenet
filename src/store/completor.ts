import { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore, type StoreDefinition } from 'pinia';
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
            return await this.listJobs().then((jobs) => jobs.find((j) => j.name === name));
        },
        async listJobs(): Promise<Job[]> {
            if (this.jobs.length > 0) return this.jobs;

            this.jobs = await this.completeJobs({});
            return this.jobs;
        },
        async completeJobs(req: CompleteJobsRequest): Promise<Job[]> {
            const { $grpc } = useNuxtApp();
            try {
                const call = $grpc.getCompletorClient().completeJobs(req);
                const { response } = await call;

                return response.jobs;
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },

        // Citizens
        async findCitizen(userId: number): Promise<UserShort | undefined> {
            return await this.completeCitizens({ userId, search: '' }).then((users) =>
                users.length === 0 ? undefined : users[0],
            );
        },
        async completeCitizens(req: CompleteCitizensRequest): Promise<UserShort[]> {
            const { $grpc } = useNuxtApp();
            try {
                const call = $grpc.getCompletorClient().completeCitizens(req);
                const { response } = await call;

                return response.users;
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },

        // Document Categories
        async completeDocumentCategories(search: string): Promise<Category[]> {
            if (!can('CompletorService.CompleteDocumentCategories')) {
                return [];
            }

            const { $grpc } = useNuxtApp();
            try {
                const call = $grpc.getCompletorClient().completeDocumentCategories({ search });
                const { response } = await call;

                return response.categories;
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },

        // Laws
        async listLawBooks(): Promise<LawBook[]> {
            const { $grpc } = useNuxtApp();
            try {
                const call = $grpc.getCompletorClient().listLawBooks({});
                const { response } = await call;

                return response.books;
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCompletorStore as unknown as StoreDefinition, import.meta.hot));
}

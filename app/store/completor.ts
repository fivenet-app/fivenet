import { defineStore } from 'pinia';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { LawBook } from '~~/gen/ts/resources/laws/laws';
import type { Job } from '~~/gen/ts/resources/users/jobs';
import type { CitizenAttribute, UserShort } from '~~/gen/ts/resources/users/users';
import type { CompleteCitizensRequest, CompleteJobsRequest } from '~~/gen/ts/services/completor/completor';
import type { ListColleaguesRequest } from '~~/gen/ts/services/jobs/jobs';

export interface CompletorState {
    jobs: Job[];
}

export const useCompletorStore = defineStore('completor', {
    state: () =>
        ({
            jobs: [],
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
            try {
                const call = getGRPCCompletorClient().completeJobs(req);
                const { response } = await call;

                return response.jobs;
            } catch (e) {
                handleGRPCError(e as RpcError);
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
            const { can } = useAuth();
            if (!can('CompletorService.CompleteCitizens').value) {
                return [];
            }

            try {
                const call = getGRPCCompletorClient().completeCitizens(req);
                const { response } = await call;

                return response.users;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Colleagues
        async findColleague(userId: number): Promise<Colleague | undefined> {
            return await this.listColleagues({
                userId: userId,
                search: '',
                labelIds: [],
            }).then((colleagues) => (colleagues.length === 0 ? undefined : colleagues[0]));
        },
        async listColleagues(req: ListColleaguesRequest): Promise<Colleague[]> {
            if (!req.pagination) {
                req.pagination = { offset: 0 };
            }

            try {
                const call = getGRPCJobsClient().listColleagues(req);
                const { response } = await call;

                return response.colleagues;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Document Categories
        async completeDocumentCategories(search: string): Promise<Category[]> {
            const { can } = useAuth();
            if (!can('CompletorService.CompleteDocumentCategories').value) {
                return [];
            }

            try {
                const call = getGRPCCompletorClient().completeDocumentCategories({
                    search: search,
                });
                const { response } = await call;

                return response.categories;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Laws
        async listLawBooks(): Promise<LawBook[]> {
            try {
                const call = getGRPCCompletorClient().listLawBooks({});
                const { response } = await call;

                return response.books;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Citizens Attributes
        async completeCitizensAttributes(search: string): Promise<CitizenAttribute[]> {
            try {
                const call = getGRPCCompletorClient().completeCitizenAttributes({
                    search: search,
                });
                const { response } = await call;

                return response.attributes;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCompletorStore, import.meta.hot));
}

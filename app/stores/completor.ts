import { defineStore } from 'pinia';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { LawBook } from '~~/gen/ts/resources/laws/laws';
import type { Label } from '~~/gen/ts/resources/users/labels';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CompleteCitizensRequest, CompleteJobsRequest } from '~~/gen/ts/services/completor/completor';
import type { ListColleaguesRequest } from '~~/gen/ts/services/jobs/jobs';

export const useCompletorStore = defineStore(
    'completor',
    () => {
        const { $grpc } = useNuxtApp();

        // State
        const jobs = ref<Job[]>([]);

        // Actions
        const getJobByName = async (name: string): Promise<Job | undefined> => {
            return listJobs().then((cachedJobs) => cachedJobs.find((j) => j.name === name));
        };

        const listJobs = async (): Promise<Job[]> => {
            if (jobs.value.length > 0) {
                return jobs.value;
            }
            jobs.value = await completeJobs({});
            return jobs.value;
        };

        const completeJobs = async (req: CompleteJobsRequest): Promise<Job[]> => {
            try {
                const call = $grpc.completor.completor.completeJobs(req);
                const { response } = await call;
                return response.jobs;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const findCitizen = async (userId: number): Promise<UserShort | undefined> => {
            const users = await completeCitizens({ userId, search: '' });
            return users.length === 0 ? undefined : users[0];
        };

        const completeCitizens = async (req: CompleteCitizensRequest): Promise<UserShort[]> => {
            const { can } = useAuth();
            if (!can('completor.CompletorService/CompleteCitizens').value) {
                return [];
            }
            try {
                const call = $grpc.completor.completor.completeCitizens(req);
                const { response } = await call;
                return response.users;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const findColleague = async (userId: number): Promise<Colleague | undefined> => {
            const colleagues = await listColleagues({ userId, search: '', labelIds: [] });
            return colleagues.length === 0 ? undefined : colleagues[0];
        };

        const listColleagues = async (req: ListColleaguesRequest): Promise<Colleague[]> => {
            if (!req.pagination) {
                req.pagination = { offset: 0 };
            }
            try {
                const call = $grpc.jobs.jobs.listColleagues(req);
                const { response } = await call;
                return response.colleagues;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const completeDocumentCategories = async (search: string): Promise<Category[]> => {
            const { can } = useAuth();
            if (!can('completor.CompletorService/CompleteDocumentCategories').value) {
                return [];
            }
            try {
                const call = $grpc.completor.completor.completeDocumentCategories({ search });
                const { response } = await call;
                return response.categories;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const listLawBooks = async (): Promise<LawBook[]> => {
            try {
                const call = $grpc.completor.completor.listLawBooks({});
                const { response } = await call;
                return response.books;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const completeCitizenLabels = async (search: string): Promise<Label[]> => {
            try {
                const call = $grpc.completor.completor.completeCitizenLabels({ search });
                const { response } = await call;
                return response.labels;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        return {
            // State
            jobs,
            // Actions
            getJobByName,
            listJobs,
            completeJobs,
            findCitizen,
            completeCitizens,
            findColleague,
            listColleagues,
            completeDocumentCategories,
            listLawBooks,
            completeCitizenLabels,
        };
    },
    {
        persist: false,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCompletorStore, import.meta.hot));
}

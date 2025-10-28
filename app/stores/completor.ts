import { defineStore } from 'pinia';
import { getCompletorCompletorClient, getJobsJobsClient } from '~~/gen/ts/clients';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { LawBook } from '~~/gen/ts/resources/laws/laws';
import type { Label } from '~~/gen/ts/resources/users/labels';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CompleteCitizensRequest, CompleteJobsRequest } from '~~/gen/ts/services/completor/completor';
import type { ListColleaguesRequest } from '~~/gen/ts/services/jobs/jobs';

/**
 * Pinia store for managing completor-related data fetching and caching.
 */
export const useCompletorStore = defineStore(
    'completor',
    () => {
        /**
         * Cached job list.
         * @type {Ref<Job[]>}
         */
        const jobs = ref<Job[]>([]);

        /**
         * Find a job by name.
         * @param {string} name - The name of the job to find.
         * @returns {Promise<Job | undefined>} - The job with the specified name, or undefined if not found.
         */
        const getJobByName = async (name: string): Promise<Job | undefined> => {
            return listJobs().then((cachedJobs) => cachedJobs.find((j) => j.name === name));
        };

        /**
         * Fetch job list and cache it in state.
         * @returns {Promise<Job[]>} - The list of jobs.
         */
        const listJobs = async (): Promise<Job[]> => {
            if (jobs.value.length > 0) {
                return jobs.value;
            }
            jobs.value = await completeJobs({});
            return jobs.value;
        };

        /**
         * Complete jobs via API.
         * @param {CompleteJobsRequest} req - The request object for completing jobs.
         * @returns {Promise<Job[]>} - The completed jobs.
         */
        const completeJobs = async (req: CompleteJobsRequest): Promise<Job[]> => {
            const completorCompletorClient = await getCompletorCompletorClient();
            try {
                const call = completorCompletorClient.completeJobs(req);
                const { response } = await call;
                return response.jobs;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Find a user by ID.
         * @param {number} userId - The ID of the user to find.
         * @returns {Promise<UserShort | undefined>} - The user with the specified ID, or undefined if not found.
         */
        const findCitizen = async (userId: number): Promise<UserShort | undefined> => {
            const users = await completeCitizens({
                search: '',
                userIds: [userId],
                userIdsOnly: true,
            });
            return users.length === 0 ? undefined : users[0];
        };

        /**
         * Complete citizens via API.
         * @param {CompleteCitizensRequest} req - The request object for completing citizens.
         * @returns {Promise<UserShort[]>} - The completed citizens.
         */
        const completeCitizens = async (req: CompleteCitizensRequest): Promise<UserShort[]> => {
            const completorCompletorClient = await getCompletorCompletorClient();
            try {
                const call = completorCompletorClient.completeCitizens(req);
                const { response } = await call;
                return response.users.map((u) => ({
                    ...u,
                    avatar: {
                        src: u.profilePicture ? `/api/filestore/${u.profilePicture}` : undefined,
                        alt: getInitials(`${u.firstname} ${u.lastname}`),
                    },
                }));
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Find a colleague by user ID.
         * @param {number} userId - The ID of the colleague to find.
         * @returns {Promise<Colleague | undefined>} - The colleague with the specified ID, or undefined if not found.
         */
        const findColleague = async (userId: number): Promise<Colleague | undefined> => {
            const colleagues = await listColleagues({
                userIds: [userId],
                search: '',
                labelIds: [],
                userOnly: true,
            });
            return colleagues.length === 0 ? undefined : colleagues[0];
        };

        /**
         * Fetch colleagues.
         * @param {ListColleaguesRequest} req - The request object for listing colleagues.
         * @returns {Promise<Colleague[]>} - The list of colleagues.
         */
        const listColleagues = async (req: ListColleaguesRequest): Promise<Colleague[]> => {
            if (!req.pagination) {
                req.pagination = { offset: 0 };
            }
            const jobsJobsClient = await getJobsJobsClient();
            try {
                const call = jobsJobsClient.listColleagues(req);
                const { response } = await call;
                return response.colleagues;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Complete colleagues via API.
         * @param {string} search - The search term for completing colleagues.
         * @param {number[]} [userIds=[]] - Optional user IDs to filter by.
         * @returns {Promise<Colleague[]>} - The completed colleagues.
         */
        const completeColleagues = async (search: string, userIds: number[] = []): Promise<Colleague[]> => {
            try {
                const colleagues = await listColleagues({
                    search: search,
                    labelIds: [],
                    userIds: userIds,
                });
                return colleagues.map((c) => ({
                    ...c,
                    avatar: {
                        src: c.profilePicture ? `/api/filestore/${c.profilePicture}` : undefined,
                        alt: getInitials(`${c.firstname} ${c.lastname}`),
                    },
                }));
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Complete document categories.
         * @param {string} search - The search term for completing document categories.
         * @returns {Promise<Category[]>} - The completed document categories.
         */
        const completeDocumentCategories = async (search: string): Promise<Category[]> => {
            const { can } = useAuth();
            if (!can('completor.CompletorService/CompleteDocumentCategories').value) {
                return [];
            }
            const completorCompletorClient = await getCompletorCompletorClient();
            try {
                const call = completorCompletorClient.completeDocumentCategories({ search });
                const { response } = await call;
                return response.categories;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Fetch law books.
         * @returns {Promise<LawBook[]>} - The list of law books.
         */
        const listLawBooks = async (): Promise<LawBook[]> => {
            const completorCompletorClient = await getCompletorCompletorClient();
            try {
                const call = completorCompletorClient.listLawBooks({});
                const { response } = await call;
                return response.books;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Complete citizen labels.
         * @param {string} search - The search term for completing citizen labels.
         * @returns {Promise<Label[]>} - The completed citizen labels.
         */
        const completeCitizenLabels = async (search: string): Promise<Label[]> => {
            const completorCompletorClient = await getCompletorCompletorClient();
            try {
                const call = completorCompletorClient.completeCitizenLabels({
                    search: search,
                });
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
            completeColleagues,
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

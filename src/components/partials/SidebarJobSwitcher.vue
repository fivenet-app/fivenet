<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { Job } from '@fivenet/gen/resources/jobs/jobs_pb';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/20/solid';
import { CompleteJobsRequest } from '@fivenet/gen/services/completor/completor_pb';
import { watchDebounced } from '@vueuse/shared';
import { SetJobRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { useNotificationsStore } from '~/store/notifications';

const { $grpc } = useNuxtApp();
const { t } = useI18n();

const authStore = useAuthStore();
const notifications = useNotificationsStore();

const { activeChar } = storeToRefs(authStore);
const { setAccessToken, setActiveChar, setJobProps } = authStore;

let entriesJobs = [] as Job[];
const filteredJobs = ref<Job[]>([]);
const queryJob = ref('');
const selectedJob = ref<undefined | Job>();

async function findJobs(): Promise<void> {
    return new Promise(async (res, rej) => {
        if (entriesJobs.length > 0) {
            return res();
        }

        const req = new CompleteJobsRequest();
        req.setSearch(queryJob.value);

        try {
            const resp = await $grpc.getCompletorClient().
                completeJobs(req, null)

            entriesJobs = resp.getJobsList();
            filteredJobs.value = entriesJobs;

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e);
        }
    });
}

async function setJob(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new SetJobRequest();
        req.setCharId(activeChar.value?.getUserId()!);
        req.setJob(selectedJob.value?.getName()!);
        const grades = selectedJob.value?.getGradesList()!;
        req.setJobGrade(grades[grades.length - 1].getGrade());

        try {
            const resp = await $grpc.getAuthClient().
                setJob(req, null);

            const promises = [
                setAccessToken(resp.getToken(), toDate(resp.getExpires()) as null | Date),
                setActiveChar(resp.getChar()!),
            ];
            if (resp.hasJobProps()) {
                promises.push(setJobProps(resp.getJobProps()!));
            } else {
                setJobProps(null);
            }
            await Promise.all(promises);

            notifications.dispatchNotification({
                title: 'notifications.job_switcher.setjob.title',
                titleI18n: true,
                content: t('notifications.job_switcher.setjob.title', [selectedJob.value?.getLabel()]),
                type: 'info'
            });

            await navigateTo({ name: 'overview' });

            queryJob.value = '';
            return res();
        } catch (e) {
            queryJob.value = '';
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watchDebounced(queryJob, async () => { filteredJobs.value = entriesJobs.filter(g => g.getLabel().toLowerCase().includes(queryJob.value.toLowerCase())) }, { debounce: 600, maxWait: 1750 });
watchDebounced(selectedJob, () => setJob());
</script>

<template>
    <Combobox as="div" v-model="selectedJob" nullable>
        <div class="relative">
            <ComboboxButton as="div">
                <ComboboxInput @click="findJobs"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    @change="queryJob = $event.target.value" :display-value="(job: any) => job ? job?.getLabel() : ''" />
            </ComboboxButton>

            <ComboboxOptions v-if="filteredJobs.length > 0"
                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                <ComboboxOption v-for="job in filteredJobs" :key="job.getName()" :value="job" as="job"
                    v-slot="{ active, selected }">
                    <li
                        :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                        <span :class="['block truncate', selected && 'font-semibold']">
                            {{ job.getLabel() }}
                        </span>

                        <span v-if="selected"
                            :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                        </span>
                    </li>
                </ComboboxOption>
            </ComboboxOptions>
        </div>
    </Combobox>
</template>

<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { Job } from '~~/gen/ts/resources/jobs/jobs';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/20/solid';
import { watchDebounced } from '@vueuse/shared';
import { SetJobRequest } from '~~/gen/ts/services/auth/auth';
import { useNotificationsStore } from '~/store/notifications';
import { RpcError } from 'grpc-web';

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

        try {
            const call = $grpc.getCompletorClient().
                completeJobs({
                    search: queryJob.value,
                });
            const { response } = await call;

            entriesJobs = response.jobs;
            filteredJobs.value = entriesJobs;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e);
        }
    });
}

async function setJob(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const grades = selectedJob.value?.grades!;

            const call = $grpc.getAuthClient().
                setJob({
                    charId: activeChar.value?.userId!,
                    job: selectedJob.value?.name!,
                    jobGrade: grades[grades.length - 1].grade,
                });
            const { response } = await call;

            const promises = [
                setAccessToken(response.token, toDate(response.expires) as null | Date),
                setActiveChar(response.char!),
            ];
            if (response.jobProps) {
                promises.push(setJobProps(response.jobProps!));
            } else {
                setJobProps(null);
            }
            await Promise.all(promises);

            notifications.dispatchNotification({
                title: 'notifications.job_switcher.setjob.title',
                titleI18n: true,
                content: t('notifications.job_switcher.setjob.title', [selectedJob.value?.label]),
                type: 'info'
            });

            await navigateTo({ name: 'overview' });

            queryJob.value = '';
            return res();
        } catch (e) {
            queryJob.value = '';
            return rej(e as RpcError);
        }
    });
}

watchDebounced(queryJob, async () => { filteredJobs.value = entriesJobs.filter(g => g.label.toLowerCase().includes(queryJob.value.toLowerCase())) }, { debounce: 600, maxWait: 1750 });
watchDebounced(selectedJob, () => setJob());
</script>

<template>
    <Combobox as="div" v-model="selectedJob" nullable>
        <div class="relative">
            <ComboboxButton as="div">
                <ComboboxInput @click="findJobs"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    @change="queryJob = $event.target.value" :display-value="(job: any) => job ? job?.label : ''" />
            </ComboboxButton>

            <ComboboxOptions v-if="filteredJobs.length > 0"
                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                <ComboboxOption v-for="job in filteredJobs" :key="job.name" :value="job" as="job"
                    v-slot="{ active, selected }">
                    <li
                        :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                        <span :class="['block truncate', selected && 'font-semibold']">
                            {{ job.label }}
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

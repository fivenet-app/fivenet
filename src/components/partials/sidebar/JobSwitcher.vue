<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { CheckIcon } from 'mdi-vue3';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificationsStore } from '~/store/notifications';
import { Job } from '~~/gen/ts/resources/users/jobs';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const notifications = useNotificationsStore();

const { activeChar } = storeToRefs(authStore);
const { setAccessToken, setActiveChar, setJobProps } = authStore;

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const queryJob = ref('');
const selectedJob = ref<undefined | Job>();

async function setJob(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const grades = selectedJob.value?.grades;
            if (!grades) return;

            const call = $grpc.getAuthClient().setJob({
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
                title: { key: 'notifications.job_switcher.setjob.title', parameters: [] },
                content: { key: 'notifications.job_switcher.setjob.title', parameters: [selectedJob.value?.label!] },
                type: 'info',
            });

            await navigateTo({ name: 'overview' });

            queryJob.value = '';
            return res();
        } catch (e) {
            queryJob.value = '';
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const filteredJobs = computed(() => jobs.value.filter((g) => g.label.toLowerCase().includes(queryJob.value.toLowerCase())));

watch(selectedJob, () => setJob());
</script>

<template>
    <Combobox as="div" v-model="selectedJob" nullable>
        <div class="relative">
            <ComboboxButton as="div">
                <ComboboxInput
                    @click="listJobs"
                    class="hidden md:block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    @change="queryJob = $event.target.value"
                    :display-value="(job: any) => (job ? job?.label : '')"
                    :placeholder="`${$t('common.select')} ${$t('common.job')}`"
                />
            </ComboboxButton>

            <ComboboxOptions
                v-if="filteredJobs.length > 0"
                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
            >
                <ComboboxOption v-for="job in filteredJobs" :key="job.name" :value="job" as="job" v-slot="{ active, selected }">
                    <li
                        :class="[
                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                            active ? 'bg-primary-500' : '',
                        ]"
                    >
                        <span :class="['block truncate', selected && 'font-semibold']">
                            {{ job.label }}
                        </span>

                        <span
                            v-if="selected"
                            :class="[
                                active ? 'text-neutral' : 'text-primary-500',
                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                            ]"
                        >
                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                        </span>
                    </li>
                </ComboboxOption>
            </ComboboxOptions>
        </div>
    </Combobox>
</template>

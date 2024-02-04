<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { CheckIcon } from 'mdi-vue3';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Job } from '~~/gen/ts/resources/users/jobs';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const notifications = useNotificatorStore();

const { activeChar } = storeToRefs(authStore);
const { setAccessToken, setActiveChar, setJobProps } = authStore;

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const queryJob = ref('');
const selectedJob = ref<undefined | Job>();

async function setJob(): Promise<void> {
    try {
        const grades = selectedJob.value?.grades;
        if (!grades) {
            return;
        }

        const call = $grpc.getAuthClient().setJob({
            charId: activeChar.value!.userId,
            job: selectedJob.value!.name,
            jobGrade: grades[grades.length - 1].grade,
        });
        const { response } = await call;

        setAccessToken(response.token, toDate(response.expires) as null | Date);
        setActiveChar(response.char!);
        setJobProps(response.jobProps);

        notifications.dispatchNotification({
            title: { key: 'notifications.job_switcher.setjob.title', parameters: {} },
            content: { key: 'notifications.job_switcher.setjob.title', parameters: { job: selectedJob.value!.label! } },
            type: 'info',
        });

        await navigateTo({ name: 'overview' });

        queryJob.value = '';
    } catch (e) {
        queryJob.value = '';
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const filteredJobs = computed(() => jobs.value.filter((g) => g.label.toLowerCase().includes(queryJob.value.toLowerCase())));

watch(selectedJob, () => setJob());
</script>

<template>
    <Combobox v-model="selectedJob" as="div" nullable>
        <div class="relative">
            <ComboboxButton as="div">
                <ComboboxInput
                    autocomplete="off"
                    class="hidden w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6 md:block"
                    :display-value="(job: any) => (job ? job?.label : '')"
                    :placeholder="`${$t('common.select')} ${$t('common.job')}`"
                    @click="listJobs"
                    @change="queryJob = $event.target.value"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </ComboboxButton>

            <ComboboxOptions
                v-if="filteredJobs.length > 0"
                class="absolute z-40 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
            >
                <ComboboxOption v-for="job in filteredJobs" :key="job.name" v-slot="{ active, selected }" :value="job">
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
                            <CheckIcon class="h-5 w-5" aria-hidden="true" />
                        </span>
                    </li>
                </ComboboxOption>
            </ComboboxOptions>
        </div>
    </Combobox>
</template>

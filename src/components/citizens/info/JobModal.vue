<script lang="ts" setup>
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Dialog,
    DialogPanel,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { watchDebounced } from '@vueuse/core';
import { CheckIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificationsStore } from '~/store/notifications';
import { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();
const notifications = useNotificationsStore();

const props = defineProps<{
    open: boolean;
    user: User;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const queryJob = ref<string>('');
const selectedJob = ref<undefined | Job>();
const selectedJobGrade = ref<undefined | JobGrade>();

const { data: jobs } = useLazyAsyncData('jobs', () => getJobs());

async function getJobs(): Promise<Array<Job>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeJobs({
                search: queryJob.value,
            });
            const { response } = await call;

            return res(response.jobs);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(jobs, () => {
    if (jobs.value) {
        selectedJob.value = jobs.value.find((j) => j.name === props.user.job);
        selectedJobGrade.value = selectedJob.value?.grades.find((g) => g.grade === props.user.jobGrade);
    }
});

watchDebounced(queryJob, async () => await getJobs(), {
    debounce: 600,
    maxWait: 1750,
});

async function setJobProp(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!selectedJob.value || selectedJob.value.name === props.user.job) return res();

        const userProps: UserProps = {
            userId: props.user.userId,
            jobName: selectedJob.value.name,
            jobGradeNumber: selectedJobGrade.value ? selectedJobGrade.value?.grade : 1,
        };

        try {
            await $grpc.getCitizenStoreClient().setUserProps({
                props: userProps,
                reason: values.reason,
            });

            props.user.job = selectedJob.value?.name!;
            props.user.jobLabel = selectedJob.value?.label!;

            props.user.jobGrade = selectedJobGrade.value?.grade!;
            props.user.jobGradeLabel = selectedJob.value?.label!;

            notifications.dispatchNotification({
                title: { key: 'notifications.action_successfull.title', parameters: [] },
                content: { key: 'notifications.action_successfull.content', parameters: [] },
                type: 'success',
            });

            emits('close');
            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

interface FormData {
    reason: string;
}

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});

const onSubmit = handleSubmit(async (values): Promise<void> => await setJobProp(values));
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <DialogPanel
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-2xl sm:p-6 h-96"
                        >
                            <form @submit="onSubmit">
                                <div class="my-2 space-y-24">
                                    <div class="flex-1 form-control">
                                        <label for="reason" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.reason') }}
                                        </label>
                                        <VeeField
                                            type="text"
                                            name="reason"
                                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.reason')"
                                            :label="$t('common.reason')"
                                        />
                                        <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="my-2">
                                    <div class="flex-1 form-control">
                                        <label for="job" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.job') }}
                                        </label>
                                        <Combobox v-if="selectedJob" as="div" v-model="selectedJob" nullable>
                                            <div class="relative">
                                                <ComboboxButton as="div">
                                                    <ComboboxInput
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @change="queryJob = $event.target.value"
                                                        :display-value="(job: any) => job.label"
                                                        autocomplete="off"
                                                    />
                                                </ComboboxButton>

                                                <ComboboxOptions
                                                    v-if="jobs"
                                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                >
                                                    <ComboboxOption
                                                        v-for="job in jobs"
                                                        :key="job.name"
                                                        :value="job"
                                                        as="char"
                                                        v-slot="{ active, selected }"
                                                    >
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
                                    </div>
                                    <div class="flex-1 form-control">
                                        <label for="job" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.job_grade') }}
                                        </label>
                                        <Combobox as="div" v-model="selectedJobGrade" nullable>
                                            <div class="relative">
                                                <ComboboxButton as="div">
                                                    <ComboboxInput
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @change="queryJob = $event.target.value"
                                                        :display-value="(grade: any) => grade.label"
                                                        autocomplete="off"
                                                    />
                                                </ComboboxButton>

                                                <ComboboxOptions
                                                    v-if="jobs"
                                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                >
                                                    <ComboboxOption
                                                        v-for="grade in selectedJob?.grades"
                                                        :key="grade.grade"
                                                        :value="grade"
                                                        as="char"
                                                        v-slot="{ active, selected }"
                                                    >
                                                        <li
                                                            :class="[
                                                                'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                active ? 'bg-primary-500' : '',
                                                            ]"
                                                        >
                                                            <span :class="['block truncate', selected && 'font-semibold']">
                                                                {{ grade.label }}
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
                                    </div>
                                </div>
                                <div class="absolute bottom-0 w-full left-0 sm:flex">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-bd bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="flex-1 rounded-bd py-2.5 px-3.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid"
                                        :class="[
                                            !meta.valid
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
                                        {{ $t('common.save') }}
                                    </button>
                                </div>
                            </form>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

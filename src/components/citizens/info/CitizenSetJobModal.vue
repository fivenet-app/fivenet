<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { max, min, required } from '@vee-validate/rules';
import { CheckIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const queryJob = ref<string>('');
const selectedJob = ref<undefined | Job>();
const selectedJobGrade = ref<undefined | JobGrade>();

watch(jobs, () => {
    selectedJob.value = jobs.value.find((j) => j.name === props.user.job);
    selectedJobGrade.value = selectedJob.value?.grades.find((g) => g.grade === props.user.jobGrade);
});

const filteredJobs = computed(() =>
    jobs.value.filter(
        (j) =>
            j.name.toLowerCase().includes(queryJob.value.toLowerCase()) ||
            j.label.toLowerCase().includes(queryJob.value.toLowerCase()),
    ),
);

const queryJobGrade = ref<string>('');

interface FormData {
    reason: string;
}

async function setJobProp(values: FormData): Promise<void> {
    if (
        selectedJob.value === undefined ||
        (selectedJob.value.name === props.user.job && selectedJobGrade.value?.grade === props.user.jobGrade)
    ) {
        isOpen.value = false;
        return;
    }

    const jobGrade = selectedJobGrade.value?.grade ? selectedJobGrade.value?.grade : 1;

    const userProps: UserProps = {
        userId: props.user.userId,
        jobName: selectedJob.value.name,
        jobGradeNumber: jobGrade,
    };

    try {
        await $grpc.getCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });

        emit('update:job', { job: selectedJob.value, grade: selectedJobGrade.value ?? { grade: 1, label: '' } });

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(selectedJob, () => {
    selectedJobGrade.value = selectedJob.value?.grades[0] ?? undefined;
});

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> => await setJobProp(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

onBeforeMount(async () => listJobs());
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.citizens.citizen_info_profile.set_job') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}" @submit.prevent="onSubmitThrottle">
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="reason" class="block text-sm font-medium leading-6">
                                {{ $t('common.reason') }}
                            </label>
                            <VeeField
                                type="text"
                                name="reason"
                                class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.reason')"
                                :label="$t('common.reason')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                    <div class="my-2">
                        <div class="flex-1">
                            <label for="job" class="block text-sm font-medium leading-6">
                                {{ $t('common.job') }}
                            </label>
                            <Combobox v-model="selectedJob" as="div" nullable>
                                <div class="relative">
                                    <ComboboxButton as="div">
                                        <ComboboxInput
                                            autocomplete="off"
                                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :display-value="(job: any) => job.label"
                                            @change="queryJob = $event.target.value"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </ComboboxButton>

                                    <ComboboxOptions
                                        v-if="filteredJobs"
                                        class="absolute z-20 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                    >
                                        <ComboboxOption
                                            v-for="job in filteredJobs"
                                            :key="job.name"
                                            v-slot="{ active, selected }"
                                            :value="job"
                                            as="char"
                                        >
                                            <li
                                                :class="[
                                                    'relative cursor-default select-none py-2 pl-8 pr-4',
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
                                                    <CheckIcon class="size-5" />
                                                </span>
                                            </li>
                                        </ComboboxOption>
                                    </ComboboxOptions>
                                </div>
                            </Combobox>
                        </div>
                        <div class="flex-1">
                            <label for="jobGrade" class="block text-sm font-medium leading-6">
                                {{ $t('common.job_grade') }}
                            </label>
                            <Combobox v-model="selectedJobGrade" as="div">
                                <div class="relative">
                                    <ComboboxButton as="div">
                                        <ComboboxInput
                                            autocomplete="off"
                                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :display-value="(grade: any) => grade?.label ?? 'N/A'"
                                            @change="queryJobGrade = $event.target.value"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </ComboboxButton>

                                    <ComboboxOptions
                                        v-if="selectedJob"
                                        class="absolute z-20 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                    >
                                        <ComboboxOption
                                            v-for="grade in selectedJob?.grades.filter((g) =>
                                                g.label.toLowerCase().includes(queryJobGrade.toLowerCase()),
                                            )"
                                            :key="grade.grade"
                                            v-slot="{ active, selected }"
                                            :value="grade"
                                            as="char"
                                        >
                                            <li
                                                :class="[
                                                    'relative cursor-default select-none py-2 pl-8 pr-4',
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
                                                    <CheckIcon class="size-5" />
                                                </span>
                                            </li>
                                        </ComboboxOption>
                                    </ComboboxOptions>
                                </div>
                            </Combobox>
                        </div>
                    </div>
                </UForm>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.save') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

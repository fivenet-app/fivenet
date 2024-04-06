<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
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

const selectedJob = ref<undefined | Job>();
const selectedJobGrade = ref<undefined | JobGrade>();

watch(jobs, () => {
    selectedJob.value = jobs.value.find((j) => j.name === props.user.job);
    selectedJobGrade.value = selectedJob.value?.grades.find((g) => g.grade === props.user.jobGrade);
});

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
                        {{ $t('components.citizens.CitizenInfoProfile.set_job') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}">
                    <UFormGroup class="flex-1" name="reason" :label="$t('common.reason')">
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
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="job" :label="$t('common.job')">
                        <USelectMenu v-model="selectedJob" :options="jobs" by="label">
                            <template #label>
                                <template v-if="selectedJob">
                                    <span class="truncate">{{ selectedJob?.label }} ({{ selectedJob.name }})</span>
                                </template>
                            </template>
                            <template #option="{ option: job }">
                                <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                            </template>
                        </USelectMenu>
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="grade" :label="$t('common.job_grade')">
                        <USelectMenu v-model="selectedJobGrade" :options="selectedJob?.grades" by="grade">
                            <template #label>
                                <template v-if="selectedJobGrade">
                                    <span class="truncate">{{ selectedJobGrade?.label }} ({{ selectedJobGrade?.grade }})</span>
                                </template>
                            </template>
                            <template #option="{ option: jobGrade }">
                                <span class="truncate">{{ jobGrade.label }} ({{ jobGrade.grade }})</span>
                            </template>
                        </USelectMenu>
                    </UFormGroup>
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

<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { useCompletorStore } from '~/stores/completor';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props/props';
import type { User } from '~~/gen/ts/resources/users/user';

const props = defineProps<{
    user: User;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
}>();

const { game } = useAppConfig();

const { t } = useI18n();

const notifications = useNotificationsStore();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const citizensCitizensClient = await getCitizensCitizensClient();

const overlay = useOverlay();
const resetConfirmModal = overlay.create(ConfirmModal);

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
    job: z.custom<Job>().optional(),
    grade: z.custom<JobGrade>().optional(),
    reset: z.coerce.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    job: undefined,
    grade: undefined,
    reset: false,
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

async function setJobProp(values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
        jobName: values.job?.name,
        jobGradeNumber: values.grade?.grade,
    };

    if (values.reset) {
        userProps.job = undefined;
        userProps.jobName = undefined;
        userProps.jobGrade = undefined;
        userProps.jobGradeNumber = undefined;
    }

    try {
        const call = citizensCitizensClient.setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emits('update:job', {
            job: response.props?.job ?? { name: '', label: '', grades: [] },
            grade: response.props?.jobGrade ?? { grade: game.startJobGrade, label: '' },
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await setJobProp(event.data);
});

const formRef = useTemplateRef('formRef');

function syncJobState(): void {
    const job = jobs.value.find((j) => j.name === props.user.job);

    state.job = job;
    state.grade = job?.grades.find((g) => g.grade === props.user.jobGrade) ?? job?.grades[0];
}

watch(
    () => state.job?.name,
    (jobName, previousJobName) => {
        if (!jobName || jobName === previousJobName) return;

        const job = jobs.value.find((j) => j.name === jobName);
        state.grade = jobName === props.user.job ? job?.grades.find((g) => g.grade === props.user.jobGrade) : job?.grades[0];
    },
);

function submitChanges(): void {
    state.reset = false;
    formRef.value?.submit();
}

function confirmReset(): void {
    resetConfirmModal.open({
        title: t('components.citizens.CitizenInfoProfile.reset_job'),
        description: t('components.citizens.CitizenInfoProfile.reset_job_description'),
        confirm: () => {
            state.reset = true;
            formRef.value?.submit();
        },
    });
}

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emits('close', false);
}

onBeforeMount(async () => {
    await listJobs();
    syncJobState();
});
</script>

<template>
    <UDrawer
        :title="$t('components.citizens.CitizenInfoProfile.set_job')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
        :ui="{ body: 'mx-auto w-full max-w-xl' }"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.citizens.CitizenInfoProfile.set_job') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="grid gap-4" :schema="schema" :state="state" @submit="submit">
                <UFormField class="flex-1" name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" class="w-full" type="text" :placeholder="$t('common.reason')" />
                </UFormField>

                <UFormField class="flex-1" name="job" :label="$t('common.job')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.job"
                            class="w-full"
                            :items="jobs"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['label', 'name']"
                        >
                            <template v-if="state.job" #default>
                                {{ state.job?.label }}
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.job')]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField class="flex-1" name="grade" :label="$t('common.job_grade')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.grade"
                            class="w-full"
                            :disabled="!state.job"
                            :items="state.job?.grades"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        >
                            <template v-if="state.grade" #default>
                                {{ state.grade?.label }} ({{ state.grade?.grade }})
                            </template>

                            <template #item-label="{ item }"> {{ item.label }} ({{ item.grade }}) </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.job_grade')]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <div class="flex w-full flex-col-reverse gap-2 sm:flex-row">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="submitChanges"
                />

                <UButton
                    class="flex-1"
                    color="error"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('components.citizens.CitizenInfoProfile.reset_job')"
                    @click="confirmReset"
                />

                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>
    </UDrawer>
</template>

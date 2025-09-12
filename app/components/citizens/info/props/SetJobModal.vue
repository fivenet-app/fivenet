<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { useCompletorStore } from '~/stores/completor';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
}>();

const { game } = useAppConfig();

const notifications = useNotificationsStore();

const completorStore = useCompletorStore();

const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    reason: z.string().min(3).max(255),
    job: z.custom<Job>().optional(),
    grade: z.custom<JobGrade>().optional(),
    reset: z.coerce.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    job: jobs.value.find((j) => j.name === props.user.job) ?? { name: '', label: '', grades: [] },
    grade: jobs.value.find((j) => j.name === props.user.job)?.grades.find((g) => g.grade === props.user.jobGrade),
    reset: false,
});

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

        emit('update:job', {
            job: response.props?.job ?? { name: '', label: '', grades: [] },
            grade: response.props?.jobGrade ?? { grade: game.startJobGrade, label: '' },
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setJobProp(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

onBeforeMount(async () => listJobs());
</script>

<template>
    <UModal :title="$t('components.citizens.CitizenInfoProfile.set_job')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                </UFormField>

                <UFormField class="flex-1" name="job" :label="$t('common.job')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.job"
                            :items="jobs"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['label', 'name']"
                        >
                            <template v-if="state.job" #default>
                                {{ state.job?.label }}
                            </template>

                            <template #item="{ item }">
                                {{ item.label }}
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
                            :items="state.job?.grades"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        >
                            <template v-if="state.grade" #default>
                                {{ state.grade?.label }} ({{ state.grade?.grade }})
                            </template>

                            <template #item="{ item }"> {{ item.label }} ({{ item.grade }}) </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.job_grade')]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />

                <UButton
                    color="error"
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.reset')"
                    @click="
                        state.reset = true;
                        formRef?.submit();
                    "
                />

                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UButtonGroup>
        </template>
    </UModal>
</template>

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
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
}>();

const { isOpen } = useOverlay();

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

        isOpen.value = false;
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

onBeforeMount(async () => listJobs());
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.citizens.CitizenInfoProfile.set_job') }}
                        </h3>

                        <UButton
                            class="-my-1"
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            @click="isOpen = false"
                        />
                    </div>
                </template>

                <div>
                    <UFormField class="flex-1" name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormField>

                    <UFormField class="flex-1" name="job" :label="$t('common.job')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.job"
                                :items="jobs"
                                by="label"
                                :searchable-placeholder="$t('common.search_field')"
                                :search-attributes="['label', 'name']"
                            >
                                <template #item-label>
                                    <template v-if="state.job">
                                        <span class="truncate">{{ state.job?.label }}</span>
                                    </template>
                                </template>

                                <template #option="{ option: job }">
                                    <span class="truncate">{{ job.label }}</span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
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
                                by="grade"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #item-label>
                                    <span v-if="state.grade" class="truncate"
                                        >{{ state.grade?.label }} ({{ state.grade?.grade }})</span
                                    >
                                </template>

                                <template #option="{ option: jobGrade }">
                                    <span class="truncate">{{ jobGrade.label }} ({{ jobGrade.grade }})</span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.job_grade')]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>

                        <UButton
                            class="flex-1"
                            type="submit"
                            block
                            color="error"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton class="flex-1" color="neutral" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

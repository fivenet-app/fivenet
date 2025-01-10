<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type { UserProps } from '~~/gen/ts/resources/users/props';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
}>();

const { isOpen } = useModal();

const { game } = useAppConfig();

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();

const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const schema = z.object({
    reason: z.string().min(3).max(255),
    job: z.custom<Job>().optional(),
    grade: z.custom<JobGrade>().optional(),
    reset: z.boolean(),
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
        const call = getGRPCCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:job', {
            job: response.props?.job ?? { name: '', label: '', grades: [] },
            grade: response.props?.jobGrade ?? { grade: game.startJobGrade, label: '' },
        });

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
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
                    <UFormGroup class="flex-1" name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="job" :label="$t('common.job')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.job"
                                :options="jobs"
                                by="label"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #label>
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
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="grade" :label="$t('common.job_grade')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.grade"
                                :options="state.job?.grades"
                                by="grade"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #label>
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
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>

                        <UButton
                            type="submit"
                            block
                            class="flex-1"
                            color="red"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

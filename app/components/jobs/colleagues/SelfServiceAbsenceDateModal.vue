<script lang="ts" setup>
import { CalendarDate, type DateValue } from '@internationalized/date';
import type { FormSubmitEvent } from '@nuxt/ui';
import { addDays, isBefore, isFuture, subDays } from 'date-fns';
import { z } from 'zod';
import InputDateRangePopover from '~/components/partials/InputDateRangePopover.vue';
import { useAuthStore } from '~/stores/auth';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { ColleagueProps } from '~~/gen/ts/resources/jobs/colleagues';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = defineProps<{
    userId: number;
    userProps?: ColleagueProps;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:absenceDates', value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void;
}>();

const notifications = useNotificationsStore();

const authStore = useAuthStore();
const { jobProps } = storeToRefs(authStore);

const jobsJobsClient = await getJobsJobsClient();

const today = new Date();
const minStart = subDays(today, jobProps.value?.settings?.absencePastDays ?? 7);
const maxEnd = addDays(today, jobProps.value?.settings?.absenceFutureDays ?? 93);

const schema = z.union([
    z.object({
        reason: z.string().min(3).max(255),
        absence: z.object({
            start: z.date().min(minStart).max(maxEnd),
            end: z.date().min(minStart).max(maxEnd),
        }),
        reset: z.literal(false),
    }),
    z.object({
        reason: z.string().min(3).max(255),
        absence: z
            .object({
                start: z.date(),
                end: z.date(),
            })
            .optional(),
        reset: z.literal(true),
    }),
]);

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    absence: {
        start: today,
        end: addDays(today, 1),
    },
    reset: false,
});

async function setAbsenceDate(values: Schema): Promise<void> {
    const userProps: ColleagueProps = {
        userId: props.userId,
        job: '',
        absenceBegin: values.absence?.start ? toTimestamp(values.absence.start) : {},
        absenceEnd: values.absence?.end ? toTimestamp(values.absence.end) : {},
    };

    if (values.reset) {
        userProps.absenceBegin = {};
        userProps.absenceEnd = {};
    }

    try {
        const call = jobsJobsClient.setColleagueProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:absenceDates', {
            userId: props.userId,
            absenceBegin: response.props?.absenceBegin,
            absenceEnd: response.props?.absenceEnd,
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

function updateAbsenceDateField(): void {
    if (
        props.userProps?.absenceBegin &&
        isFuture(toDate(props.userProps.absenceBegin)) &&
        props.userProps?.absenceEnd &&
        isFuture(toDate(props.userProps.absenceEnd))
    ) {
        if (!state.absence)
            state.absence = {
                start: toDate(props.userProps.absenceBegin),
                end: toDate(props.userProps.absenceEnd),
            };
    } else {
        state.absence = undefined;
    }
}

const isDateDisabled = (date: DateValue) => {
    return isBefore(date.toDate('UTC'), subDays(today, 1));
};

watch(props, () => updateAbsenceDateField());

updateAbsenceDateField();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setAbsenceDate(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.jobs.self_service.set_absence_date')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
                <UFormField name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                </UFormField>

                <UFormField name="absenceBegin" :label="$t('common.time_range')">
                    <InputDateRangePopover
                        v-model="state.absence"
                        range
                        :is-date-disabled="isDateDisabled"
                        :min-value="new CalendarDate(today.getFullYear(), today.getMonth() + 1, today.getDate())"
                        :max-value="new CalendarDate(maxEnd.getFullYear(), maxEnd.getMonth() + 1, maxEnd.getDate())"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    color="error"
                    class="flex-1"
                    block
                    :disabled="!canSubmit || (!userProps?.absenceBegin && !userProps?.absenceEnd)"
                    :loading="!canSubmit"
                    :label="$t('common.annul')"
                    @click="
                        state.reset = true;
                        formRef?.submit();
                    "
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>

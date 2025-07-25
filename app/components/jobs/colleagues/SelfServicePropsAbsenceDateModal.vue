<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { addDays, isFuture, subDays } from 'date-fns';
import { z } from 'zod';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import { useAuthStore } from '~/stores/auth';
import type { ColleagueProps } from '~~/gen/ts/resources/jobs/colleagues';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = defineProps<{
    userId: number;
    userProps?: ColleagueProps;
}>();

const emit = defineEmits<{
    (e: 'update:absenceDates', value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificationsStore();

const authStore = useAuthStore();
const { jobProps } = storeToRefs(authStore);

const today = new Date();
const minStart = subDays(today, jobProps.value?.settings?.absencePastDays ?? 7);
const maxEnd = addDays(today, jobProps.value?.settings?.absenceFutureDays ?? 93);

const schema = z.union([
    z.object({
        reason: z.string().min(3).max(255),
        absenceBegin: z.date().min(minStart),
        absenceEnd: z.date().min(today).max(maxEnd),
        reset: z.literal(false),
    }),
    z.object({
        reason: z.string().min(3).max(255),
        absenceBegin: z.date().optional(),
        absenceEnd: z.date().optional(),
        reset: z.literal(true),
    }),
]);

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    absenceBegin: today,
    absenceEnd: addDays(today, 1),
    reset: false,
});

async function setAbsenceDate(values: Schema): Promise<void> {
    const userProps: ColleagueProps = {
        userId: props.userId,
        job: '',
        absenceBegin: values.absenceBegin ? toTimestamp(values.absenceBegin) : {},
        absenceEnd: values.absenceEnd ? toTimestamp(values.absenceEnd) : {},
    };

    if (values.reset) {
        userProps.absenceBegin = {};
        userProps.absenceEnd = {};
    }

    try {
        const call = $grpc.jobs.jobs.setColleagueProps({
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

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function updateAbsenceDateField(): void {
    if (props.userProps?.absenceBegin && isFuture(toDate(props.userProps.absenceBegin))) {
        state.absenceBegin = toDate(props.userProps.absenceBegin);
    }

    if (props.userProps?.absenceEnd && isFuture(toDate(props.userProps.absenceEnd))) {
        state.absenceEnd = toDate(props.userProps.absenceEnd);
    }
}

watch(props, () => updateAbsenceDateField());

updateAbsenceDateField();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setAbsenceDate(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.jobs.self_service.set_absence_date') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormGroup>

                    <div class="flex flex-col gap-1 sm:flex-row">
                        <UFormGroup class="flex-1" name="absenceBegin" :label="$t('common.from')">
                            <PartialsDatePickerPopover
                                v-model="state.absenceBegin"
                                :popover="{ popper: { placement: 'bottom-start' } }"
                                :date-picker="{
                                    disabledDates: [
                                        { start: null, end: minStart },
                                        { start: maxEnd, end: null },
                                    ],
                                }"
                            />
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="absenceEnd" :label="$t('common.to')">
                            <DatePickerPopoverClient
                                v-model="state.absenceEnd"
                                :popover="{ popper: { placement: 'bottom-start' } }"
                                :date-picker="{
                                    disabledDates: [
                                        { start: null, end: subDays(today, 1) },
                                        { start: maxEnd, end: null },
                                    ],
                                }"
                            />
                        </UFormGroup>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            class="flex-1"
                            type="submit"
                            block
                            color="error"
                            :disabled="!canSubmit || (!userProps?.absenceBegin && !userProps?.absenceEnd)"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.annul') }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

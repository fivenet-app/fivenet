<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { addDays, format, isFuture } from 'date-fns';
import { useNotificatorStore } from '~/store/notificator';
import type { JobsUserProps } from '~~/gen/ts/resources/jobs/colleagues';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import DatePickerClient from '~/components/partials/DatePicker.client.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    userId: number;
    userProps?: JobsUserProps;
}>();

const emit = defineEmits<{
    (e: 'update:absenceDates', value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void;
}>();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const now = new Date();

const schema = z.union([
    z.object({
        reason: z.string().min(3).max(255),
        absenceBegin: z.date().min(now),
        absenceEnd: z.date().min(addDays(now, 1)),
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
    absenceBegin: now,
    absenceEnd: addDays(now, 1),
    reset: false,
});

async function setAbsenceDate(values: Schema): Promise<void> {
    const userProps: JobsUserProps = {
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
        const call = getGRPCJobsClient().setJobsUserProps({
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

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="reason" :label="$t('common.reason')">
                        <UInput
                            v-model="state.reason"
                            type="text"
                            :placeholder="$t('common.reason')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <div class="flex flex-col gap-1 sm:flex-row">
                        <UFormGroup class="flex-1" name="absenceBegin" :label="$t('common.from')">
                            <UPopover :popper="{ placement: 'bottom-start' }">
                                <UButton
                                    variant="outline"
                                    color="gray"
                                    block
                                    icon="i-mdi-calendar-month"
                                    :label="state.absenceBegin ? format(state.absenceBegin, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                                />

                                <template #panel="{ close }">
                                    <DatePickerClient v-model="state.absenceBegin" @close="close" />
                                </template>
                            </UPopover>
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="absenceEnd" :label="$t('common.to')">
                            <UPopover :popper="{ placement: 'bottom-start' }">
                                <UButton
                                    variant="outline"
                                    color="gray"
                                    block
                                    icon="i-mdi-calendar-month"
                                    :label="state.absenceEnd ? format(state.absenceEnd, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                                />

                                <template #panel="{ close }">
                                    <DatePickerClient v-model="state.absenceEnd" @close="close" />
                                </template>
                            </UPopover>
                        </UFormGroup>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            type="submit"
                            block
                            color="red"
                            class="flex-1"
                            :disabled="!canSubmit || (!userProps?.absenceBegin && !userProps?.absenceEnd)"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.annul') }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

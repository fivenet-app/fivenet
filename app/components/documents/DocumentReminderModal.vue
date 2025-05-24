<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { subDays } from 'date-fns';
import { z } from 'zod';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { SetDocumentReminderResponse } from '~~/gen/ts/services/documents/documents';

const props = defineProps<{
    documentId: number;
    reminderTime?: Timestamp;
}>();

const emit = defineEmits<{
    (e: 'update:reminderTime', reminderTime?: Timestamp): void;
}>();

const reminderTime = useVModel(props, 'reminderTime', emit);

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const schema = z.object({
    message: z.string().min(1).max(64),
    reminderTime: z.date().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    message: '',
    reminderTime: reminderTime.value ? toDate(reminderTime.value) : undefined,
});

watch(reminderTime, () => (state.reminderTime = reminderTime.value ? toDate(reminderTime.value) : undefined));

async function setDocumentReminder(values: Schema): Promise<SetDocumentReminderResponse> {
    try {
        const call = $grpc.documents.documents.setDocumentReminder({
            documentId: props.documentId,
            reminderTime: values.reminderTime ? toTimestamp(values.reminderTime) : undefined,
            message: values.message,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        reminderTime.value = values.reminderTime ? toTimestamp(values.reminderTime) : undefined;

        isOpen.value = false;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setDocumentReminder(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.reminder') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup
                        class="grid items-center gap-2"
                        name="reminderTime"
                        :label="$t('common.time')"
                        :ui="{ container: '' }"
                    >
                        <DatePickerPopoverClient
                            v-model="state.reminderTime"
                            date-format="dd.MM.yyyy HH:mm"
                            :popover="{ popper: { placement: 'bottom-start' } }"
                            :date-picker="{
                                mode: 'dateTime',
                                is24hr: true,
                                clearable: true,
                                disabledDates: [{ start: null, end: subDays(new Date(), 1) }],
                            }"
                        />
                    </UFormGroup>

                    <UFormGroup
                        class="grid items-center gap-2"
                        name="message"
                        :label="$t('common.message')"
                        :ui="{ container: '' }"
                    >
                        <UInput v-model="state.message" type="text" :placeholder="$t('common.message')" />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
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

<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { subDays } from 'date-fns';
import { z } from 'zod';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { SetDocumentReminderResponse } from '~~/gen/ts/services/documents/documents';

const props = defineProps<{
    documentId: number;
    reminderTime?: Timestamp;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:reminderTime', reminderTime?: Timestamp): void;
}>();

const reminderTime = useVModel(props, 'reminderTime', emit);

const notifications = useNotificationsStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const schema = z.object({
    message: z.string().min(1).max(64),
    reminderTime: z.coerce.date().optional(),
    maxReminderCount: z.coerce.number().int().min(1).max(10).default(10),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    message: '',
    reminderTime: reminderTime.value ? toDate(reminderTime.value) : undefined,
    maxReminderCount: 1,
});

watch(reminderTime, () => (state.reminderTime = reminderTime.value ? toDate(reminderTime.value) : undefined));

async function setDocumentReminder(values: Schema): Promise<SetDocumentReminderResponse> {
    try {
        const call = documentsDocumentsClient.setDocumentReminder({
            documentId: props.documentId,
            reminderTime: values.reminderTime ? toTimestamp(values.reminderTime) : undefined,
            message: values.message,
            maxReminderCount: values.maxReminderCount,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        reminderTime.value = values.reminderTime ? toTimestamp(values.reminderTime) : undefined;

        emit('close', false);

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

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('common.reminder', 2)">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField
                    class="grid items-center gap-2"
                    name="reminderTime"
                    :label="$t('common.time')"
                    :ui="{ container: '' }"
                >
                    <DatePickerPopoverClient
                        v-model="state.reminderTime"
                        date-format="dd.MM.yyyy HH:mm"
                        :date-picker="{
                            mode: 'dateTime',
                            is24hr: true,
                            clearable: true,
                            disabledDates: [{ start: null, end: subDays(new Date(), 1) }],
                        }"
                    />
                </UFormField>

                <UFormField
                    class="grid items-center gap-2"
                    name="message"
                    :label="$t('common.message')"
                    :ui="{ container: '' }"
                >
                    <UInput v-model="state.message" type="text" :placeholder="$t('common.message')" />
                </UFormField>

                <!--
                    Only show if recurring reminders are enabled
                    <UFormField
                        class="grid items-center gap-2"
                        name="message"
                        label="Max number of total reminders"
                        :ui="{ container: '' }"
                    >
                        <UInput v-model="state.maxReminderCount" type="number" :min="1" :max="10" :step="1" />
                    </UFormField>
                    -->
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>

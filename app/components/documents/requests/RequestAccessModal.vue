<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    documentId: number;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const notifications = useNotificationsStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const accessLevels = [
    { label: t(`enums.documents.AccessLevel.${AccessLevel[AccessLevel.VIEW]}`), value: AccessLevel.VIEW },
    { label: t(`enums.documents.AccessLevel.${AccessLevel[AccessLevel.COMMENT]}`), value: AccessLevel.COMMENT },
    { label: t(`enums.documents.AccessLevel.${AccessLevel[AccessLevel.STATUS]}`), value: AccessLevel.STATUS },
    { label: t(`enums.documents.AccessLevel.${AccessLevel[AccessLevel.ACCESS]}`), value: AccessLevel.ACCESS },
    { label: t(`enums.documents.AccessLevel.${AccessLevel[AccessLevel.EDIT]}`), value: AccessLevel.EDIT },
];

const schema = z.object({
    reason: z.string().min(3).max(255),
    accessLevel: z.nativeEnum(AccessLevel),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    accessLevel: AccessLevel.VIEW,
});

async function createDocumentRequest(values: Schema): Promise<void> {
    try {
        const call = documentsDocumentsClient.createDocumentReq({
            documentId: props.documentId,
            requestType: DocActivityType.REQUESTED_ACCESS,
            reason: values.reason,
            data: {
                data: {
                    oneofKind: 'accessRequested',
                    accessRequested: {
                        level: values.accessLevel,
                    },
                },
            },
        });
        await call;

        notifications.add({
            title: { key: 'notifications.documents.requests.created.title' },
            description: { key: 'notifications.documents.requests.created.content' },
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
    await createDocumentRequest(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('common.request')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                </UFormField>

                <UFormField class="flex-1" name="requestType" :label="$t('common.access')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.accessLevel"
                            :items="accessLevels"
                            :placeholder="$t('common.access')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            value-key="value"
                            class="w-full"
                        >
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
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
                    :label="$t('common.request', 2)"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>

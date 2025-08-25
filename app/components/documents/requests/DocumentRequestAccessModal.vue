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

const { isOpen } = useOverlay();

const notifications = useNotificationsStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const accessLevels = [
    { level: AccessLevel.VIEW },
    { level: AccessLevel.COMMENT },
    { level: AccessLevel.STATUS },
    { level: AccessLevel.ACCESS },
    { level: AccessLevel.EDIT },
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

        isOpen.value = false;
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
</script>

<template>
    <UModal>
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl leading-6 font-semibold">
                            {{ $t('common.request') }}
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
                    <UFormField name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormField>

                    <UFormField class="flex-1" name="requestType" :label="$t('common.access')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.accessLevel"
                                :items="accessLevels"
                                :placeholder="$t('common.access')"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #item-label>
                                    <span v-if="state.accessLevel" class="truncate">{{
                                        $t(`enums.documents.AccessLevel.${AccessLevel[state.accessLevel]}`)
                                    }}</span>
                                </template>

                                <template #item="{ option }">
                                    <span class="truncate">{{ $t(`enums.documents.AccessLevel.${AccessLevel[option]}`) }}</span>
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="neutral" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.request', 2) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    documentId: number;
}>();

const { isOpen } = useModal();

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
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.request') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="requestType" :label="$t('common.access')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.accessLevel"
                                :options="accessLevels"
                                :placeholder="$t('common.access')"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #label>
                                    <span v-if="state.accessLevel" class="truncate">{{
                                        $t(`enums.documents.AccessLevel.${AccessLevel[state.accessLevel]}`)
                                    }}</span>
                                </template>

                                <template #option="{ option }">
                                    <span class="truncate">{{ $t(`enums.documents.AccessLevel.${AccessLevel[option]}`) }}</span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
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

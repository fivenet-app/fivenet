<script lang="ts" setup>
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { DocRequest } from '~~/gen/ts/resources/documents/requests';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    request: DocRequest;
    canUpdate: boolean;
    canDelete: boolean;
}>();

const emits = defineEmits<{
    (e: 'refreshRequests'): void;
}>();

const notifications = useNotificatorStore();

async function updateDocumentReq(documentId: string, requestId: string, accepted: boolean): Promise<void> {
    try {
        const call = getGRPCDocStoreClient().updateDocumentReq({
            documentId,
            requestId,
            accepted,
        });
        const { response } = await call;

        emits('refreshRequests');

        if (response.request !== undefined) {
            if (response.request.requestType === DocActivityType.REQUESTED_UPDATE) {
                navigateTo({ name: 'documents-id-edit', params: { id: documentId } });
            }
        }

        notifications.add({
            title: { key: 'notifications.docstore.requests.updated.title' },
            description: { key: 'notifications.docstore.requests.updated.content' },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteDocumentReq(id: string): Promise<void> {
    try {
        const call = getGRPCDocStoreClient().deleteDocumentReq({
            requestId: id,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.docstore.requests.deleted.title' },
            description: { key: 'notifications.docstore.requests.deleted.content' },
            type: NotificationType.SUCCESS,
        });

        emits('refreshRequests');
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (accepted: boolean) => {
    canSubmit.value = false;
    await updateDocumentReq(props.request.documentId, props.request.id, accepted).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <li :key="request.id" class="hover:bg-neutral/5 flex justify-between gap-x-6 py-5 transition-colors">
        <div class="flex min-w-0 gap-x-2 px-2">
            <div class="min-w-0 flex-auto">
                <p class="text-base font-semibold leading-6 text-gray-100" :title="`${$t('common.id')}: ${request.id}`">
                    {{ $t(`enums.docstore.DocActivityType.${DocActivityType[request.requestType]}`) }}
                </p>
                <p class="mt-1 flex gap-1 text-sm leading-5">
                    <span class="font-semibold">{{ $t('common.reason') }}:</span> <span>{{ request.reason }}</span>
                </p>
                <p v-if="request.accepted !== undefined" class="mt-1 flex gap-1 text-sm leading-5">
                    <span class="font-semibold">{{ $t('common.accept', 2) }}:</span>
                    <span v-if="request.accepted" class="text-success-400">
                        {{ $t('common.yes') }}
                    </span>
                    <span v-else class="text-error-400">
                        {{ $t('common.no') }}
                    </span>
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-6 px-2">
            <div class="hidden text-sm sm:flex sm:flex-col sm:items-end">
                <div class="inline-flex gap-1">
                    {{ $t('common.creator') }}
                    <CitizenInfoPopover :user="request.creator" text-class="underline" />
                </div>
                <div>
                    {{ $t('common.created') }}
                    <GenericTime :value="request.createdAt" :ago="true" />
                </div>
                <div v-if="request.updatedAt">
                    {{ $t('common.updated') }}
                    <GenericTime :value="request.updatedAt" :ago="true" />
                </div>
            </div>
            <div class="flex items-center gap-2">
                <UButtonGroup v-if="canUpdate && request.accepted === undefined" class="inline-flex w-full">
                    <UButton
                        class="flex-1"
                        block
                        color="green"
                        icon="i-mdi-check-bold"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle(true)"
                    />

                    <UButton
                        class="flex-1"
                        block
                        color="red"
                        icon="i-mdi-close-thick"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle(false)"
                    />
                </UButtonGroup>

                <UDropdown
                    v-if="canDelete"
                    :items="[
                        [
                            {
                                label: $t('common.delete'),
                                icon: 'i-mdi-trash-can',
                                click: async () => deleteDocumentReq(request.id),
                            },
                        ],
                    ]"
                    :popper="{ placement: 'bottom-start' }"
                >
                    <UButton size="md" color="white" icon="i-mdi-menu" trailing-icon="i-mdi-chevron-down" />
                </UDropdown>
            </div>
        </div>
    </li>
</template>

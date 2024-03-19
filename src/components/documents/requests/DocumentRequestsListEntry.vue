<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';
import { useConfirmDialog, useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CheckBoldIcon, CloseThickIcon, MenuIcon, TrashCanIcon } from 'mdi-vue3';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { DocRequest } from '~~/gen/ts/resources/documents/requests';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';

const props = defineProps<{
    request: DocRequest;
    canUpdate: boolean;
    canDelete: boolean;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

async function updateDocumentReq(documentId: string, requestId: string, accepted: boolean): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().updateDocumentReq({
            documentId,
            requestId,
            accepted,
        });
        const { response } = await call;

        emits('refresh');

        if (response.request !== undefined) {
            if (response.request.requestType === DocActivityType.REQUESTED_UPDATE) {
                navigateTo({ name: 'documents-id-edit', params: { id: documentId } });
            }
        }

        notifications.dispatchNotification({
            title: { key: 'notifications.docstore.requests.updated.title' },
            content: { key: 'notifications.docstore.requests.updated.content' },
            type: 'success',
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteDocumentReq(id: string): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().deleteDocumentReq({
            requestId: id,
        });
        await call;

        notifications.dispatchNotification({
            title: { key: 'notifications.docstore.requests.deleted.title' },
            content: { key: 'notifications.docstore.requests.deleted.content' },
            type: 'success',
        });

        emits('refresh');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteDocumentReq(id));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (accepted: boolean) => {
    canSubmit.value = false;
    await updateDocumentReq(props.request.documentId, props.request.id, accepted).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <li :key="request.id" class="flex justify-between gap-x-6 py-5 transition-colors hover:bg-neutral/5">
        <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(request.id)" />

        <div class="flex min-w-0 gap-x-4 px-2">
            <div class="min-w-0 flex-auto">
                <p class="text-base font-semibold leading-6 text-gray-100" :title="`${$t('common.id')}: ${request.id}`">
                    {{ $t(`enums.docstore.DocActivityType.${DocActivityType[request.requestType]}`) }}
                </p>
                <p class="mt-1 flex gap-1 text-sm leading-5 text-gray-300">
                    <span class="font-semibold">{{ $t('common.reason') }}:</span> <span>{{ request.reason }}</span>
                </p>
                <p v-if="request.accepted !== undefined" class="mt-1 flex gap-1 text-sm leading-5 text-gray-300">
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
                <template v-if="canUpdate && request.accepted === undefined">
                    <button
                        type="button"
                        :disabled="!canSubmit"
                        class="rounded flex flex-1 justify-center px-3.5 py-2.5 text-sm font-semibold text-neutral"
                        :class="[
                            !canSubmit
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                        ]"
                        @click="onSubmitThrottle(true)"
                    >
                        <CheckBoldIcon class="h-5 w-5 text-success-400" aria-hidden="true" />
                    </button>
                    <button
                        type="button"
                        :disabled="!canSubmit"
                        class="rounded flex flex-1 justify-center px-2.5 py-2.5 text-sm font-semibold text-neutral"
                        :class="[
                            !canSubmit
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                        ]"
                        @click="onSubmitThrottle(false)"
                    >
                        <CloseThickIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                    </button>
                </template>

                <Menu v-if="canDelete" as="div" class="relative flex-none">
                    <MenuButton class="block text-gray-300 hover:text-gray-100">
                        <span class="sr-only">{{ $t('common.open') }}</span>
                        <MenuIcon class="h-5 w-5" aria-hidden="true" />
                    </MenuButton>
                    <transition
                        enter-active-class="transition ease-out duration-100"
                        enter-from-class="transform opacity-0 scale-95"
                        enter-to-class="transform opacity-100 scale-100"
                        leave-active-class="transition ease-in duration-75"
                        leave-from-class="transform opacity-100 scale-100"
                        leave-to-class="transform opacity-0 scale-95"
                    >
                        <MenuItems
                            class="absolute right-0 z-30 mt-2 w-28 origin-top-right rounded-md bg-base-800 py-1 shadow-float ring-1 ring-base-100 ring-opacity-5 focus:outline-none"
                        >
                            <MenuItem v-slot="{ close }">
                                <button
                                    type="button"
                                    class="inline-flex items-center px-4 py-2 text-sm text-neutral hover:transition-colors"
                                    :disabled="!canSubmit"
                                    :class="[
                                        !canSubmit
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'hover:bg-primary-400 focus-visible:outline-primary-500',
                                    ]"
                                    @click="
                                        close();
                                        reveal();
                                    "
                                >
                                    <TrashCanIcon class="h-5 w-5" aria-hidden="true" />
                                    {{ $t('common.delete') }}
                                </button>
                            </MenuItem>
                        </MenuItems>
                    </transition>
                </Menu>
            </div>
        </div>
    </li>
</template>

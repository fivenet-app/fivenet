<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';
import { CheckBoldIcon, CloseThickIcon, FrequentlyAskedQuestionsIcon, MenuIcon, TrashCanIcon } from 'mdi-vue3';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import type { ListDocumentReqsResponse } from '~~/gen/ts/services/docstore/docstore';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useNotificatorStore } from '~/store/notificator';

const props = defineProps<{
    documentId: string;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const offset = ref(0n);

const {
    data: requests,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-requests-${offset.value}`, () => listDocumnetReqs(props.documentId));

async function listDocumnetReqs(documentId: string): Promise<ListDocumentReqsResponse> {
    try {
        const call = $grpc.getDocStoreClient().listDocumentReqs({
            pagination: {
                offset: offset.value,
            },
            documentId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

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

        refresh();

        notifications.dispatchNotification({
            title: { key: 'notifications.docstore.requests.deleted.title' },
            content: { key: 'notifications.docstore.requests.deleted.content' },
            type: 'success',
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.request', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.request', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-else-if="requests === null || requests.requests.length === 0"
            :icon="FrequentlyAskedQuestionsIcon"
            :message="$t('common.not_found', [$t('common.request', 2)])"
        />

        <div v-else>
            <ul role="list" class="mb-6 divide-y divide-gray-100 rounded-md">
                <li
                    v-for="request in requests.requests"
                    :key="request.id"
                    class="flex justify-between gap-x-6 py-5 transition-colors hover:bg-neutral/5"
                >
                    <div class="flex min-w-0 gap-x-4 px-2">
                        <div class="min-w-0 flex-auto">
                            <p
                                class="text-base font-semibold leading-6 text-gray-100"
                                :title="`${$t('common.id')}: ${request.id}`"
                            >
                                {{ $t(`enums.docstore.DocActivityType.${DocActivityType[request.requestType]}`) }}
                            </p>
                            <p class="mt-1 flex text-sm leading-5 text-gray-300">
                                <span class="font-semibold">{{ $t('common.reason') }}:</span> {{ request.reason }}
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
                            <template v-if="request.accepted === undefined">
                                <button type="button" @click="updateDocumentReq(documentId, request.id, true)">
                                    <CheckBoldIcon class="h-5 w-5 text-success-400" />
                                </button>
                                <button type="button" @click="updateDocumentReq(documentId, request.id, false)">
                                    <CloseThickIcon class="h-5 w-5 text-error-400" />
                                </button>
                            </template>

                            <Menu as="div" class="relative flex-none">
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
                                                v-if="can('DocStoreService.DeleteDocumentReq')"
                                                type="button"
                                                class="inline-flex items-center"
                                                @click="
                                                    close();
                                                    deleteDocumentReq(request.id);
                                                "
                                            >
                                                <TrashCanIcon class="h-5 w-5" />
                                                {{ $t('common.delete') }}
                                            </button>
                                        </MenuItem>
                                    </MenuItems>
                                </transition>
                            </Menu>
                        </div>
                    </div>
                </li>
            </ul>
        </div>
    </div>
</template>

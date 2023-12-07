<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { MenuIcon, TrashCanIcon } from 'mdi-vue3';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import type { ListDocumentReqsResponse } from '~~/gen/ts/services/docstore/docstore';
import Time from '~/components/partials/elements/Time.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

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

async function deleteDocumentReq(id: string): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().deleteDocumentReq({
            requestId: id,
        });
        await call;

        await refresh();
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div v-if="requests !== null">
        <ul role="list" class="divide-y divide-gray-100 transition-colors hover:bg-neutral/5 rounded-md">
            <li v-for="request in requests.requests" :key="request.id" class="flex justify-between gap-x-6 py-5">
                <div class="flex min-w-0 gap-x-4 px-1">
                    <div class="min-w-0 flex-auto">
                        <p class="text-base font-semibold leading-6 text-gray-100">
                            {{ $t(`enums.docstore.DocActivityType.${DocActivityType[request.requestType]}`) }}
                        </p>
                        <p class="mt-1 flex text-xs leading-5 text-gray-300">{{ $t('common.reason') }}: {{ request.reason }}</p>
                    </div>
                </div>
                <div class="flex shrink-0 items-center gap-x-6">
                    <div class="hidden sm:flex sm:flex-col sm:items-end">
                        <CitizenInfoPopover :user="request.creator" />
                        <Time :value="request.createdAt" :ago="true" />
                    </div>
                    <Menu as="div" class="relative flex-none">
                        <MenuButton class="-m-2.5 block p-2.5 text-gray-300 hover:text-gray-100">
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
                                class="absolute z-30 right-0 w-48 py-1 mt-2 origin-top-right rounded-md shadow-float bg-base-800 ring-1 ring-base-100 ring-opacity-5 focus:outline-none"
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
                                        <TrashCanIcon class="w-6 h-6" />
                                        {{ $t('common.delete') }}
                                    </button>
                                </MenuItem>
                            </MenuItems>
                        </transition>
                    </Menu>
                </div>
            </li>
        </ul>
    </div>
</template>

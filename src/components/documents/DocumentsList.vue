<script lang="ts" setup>
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Disclosure,
    DisclosureButton,
    DisclosurePanel,
} from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCheck, mdiChevronDown } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/shared';
import { ref } from 'vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import TablePagination from '~/components/partials/TablePagination.vue';
import * as google_protobuf_timestamp_pb from '~~/gen/ts/google/protobuf/timestamp';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { DocumentCategory } from '~~/gen/ts/resources/documents/category';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ListDocumentsRequest } from '~~/gen/ts/services/docstore/docstore';
import DataNoDataBlock from '../partials/DataNoDataBlock.vue';
import DocumentsListEntry from './DocumentsListEntry.vue';
import TemplatesModal from './templates/TemplatesModal.vue';

const { $grpc } = useNuxtApp();

const search = ref<{ title: string; category?: DocumentCategory; character?: UserShort; from?: string; to?: string }>({
    title: '',
});
const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const entriesCategories = ref<DocumentCategory[]>([]);
const queryCategories = ref<string>('');
const entriesChars = ref<UserShort[]>([]);
const queryChars = ref<string>('');

const { data: documents, pending, refresh, error } = useLazyAsyncData(`documents-${offset.value}`, () => listDocuments());

async function listDocuments(): Promise<Array<DocumentShort>> {
    return new Promise(async (res, rej) => {
        const req: ListDocumentsRequest = {
            pagination: {
                offset: offset.value,
            },
            orderBy: [],
            search: search.value.title,
            categoryIds: [],
            creatorIds: [],
        };
        if (search.value.category) req.categoryIds.push(search.value.category.id);
        if (search.value.character) req.creatorIds.push(search.value.character.userId);
        if (search.value.from) {
            req.from = {
                timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(search.value.from)!),
            };
        }
        if (search.value.to) {
            req.to = {
                timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(search.value.to)!),
            };
        }

        try {
            const call = $grpc.getDocStoreClient().listDocuments(req);
            const { response } = await call;

            pagination.value = response.pagination;
            return res(response.documents);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function findCategories(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeDocumentCategories({
                search: queryCategories.value,
            });
            const { response } = await call;

            entriesCategories.value = response.categories;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function findChars(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeCitizens({
                search: queryChars.value,
            });
            const { response } = await call;

            entriesChars.value = response.users;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

const templatesOpen = ref(false);

watch(offset, async () => refresh());
watchDebounced(search.value, async () => refresh(), { debounce: 600, maxWait: 1400 });
watchDebounced(queryCategories, async () => findCategories(), {
    debounce: 600,
    maxWait: 1400,
});
watchDebounced(queryChars, async () => findChars(), {
    debounce: 600,
    maxWait: 1400,
});

onMounted(async () => {
    findCategories();
    findChars();
});
</script>

<template>
    <TemplatesModal :open="templatesOpen" @close="templatesOpen = false" />
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <label for="search" class="block mb-2 text-sm font-medium leading-6 text-neutral">
                            {{ $t('common.search') }}
                        </label>
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-1 form-control">
                                <input
                                    v-model="search.title"
                                    ref="searchInput"
                                    type="text"
                                    name="search"
                                    :placeholder="$t('common.title')"
                                    class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                />
                            </div>
                            <div class="flex-initial form-control" v-can="'DocStoreService.CreateDocument'">
                                <button
                                    @click="templatesOpen = true"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                >
                                    {{ $t('common.create') }}
                                </button>
                            </div>
                            <div class="flex-initial" v-can="'CompletorService.CompleteDocumentCategories'">
                                <NuxtLink
                                    :to="{ name: 'documents-categories' }"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                >
                                    {{ $t('common.category', 2) }}
                                </NuxtLink>
                            </div>
                            <div class="flex-initial" v-can="'DocStoreService.ListTemplates'">
                                <NuxtLink
                                    :to="{ name: 'documents-templates' }"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                >
                                    {{ $t('common.template', 2) }}
                                </NuxtLink>
                            </div>
                        </div>
                        <Disclosure as="div" class="pt-2" v-slot="{ open }">
                            <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                <span class="text-base-200 leading-7">{{ $t('common.advanced_search') }}</span>
                                <span class="ml-6 flex h-7 items-center">
                                    <SvgIcon
                                        :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                        aria-hidden="true"
                                        type="mdi"
                                        :path="mdiChevronDown"
                                    />
                                </span>
                            </DisclosureButton>
                            <DisclosurePanel class="mt-2 pr-12">
                                <div class="flex flex-row gap-2">
                                    <div class="flex-1 form-control">
                                        <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.category', 1) }}
                                        </label>
                                        <Combobox as="div" v-model="search.category" class="mt-2" nullable>
                                            <div class="relative">
                                                <ComboboxButton as="div">
                                                    <ComboboxInput
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @change="queryCategories = $event.target.value"
                                                        :display-value="(category: any) => category?.name"
                                                        :placeholder="$t('common.category', 1)"
                                                    />
                                                </ComboboxButton>

                                                <ComboboxOptions
                                                    v-if="entriesCategories.length > 0"
                                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm"
                                                >
                                                    <ComboboxOption
                                                        v-for="category in entriesCategories"
                                                        :key="category.id?.toString()"
                                                        :value="category"
                                                        as="category"
                                                        v-slot="{ active, selected }"
                                                    >
                                                        <li
                                                            :class="[
                                                                'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                active ? 'bg-primary-500' : '',
                                                            ]"
                                                        >
                                                            <span :class="['block truncate', selected && 'font-semibold']">
                                                                {{ category.name }}
                                                            </span>

                                                            <span
                                                                v-if="selected"
                                                                :class="[
                                                                    active ? 'text-neutral' : 'text-primary-500',
                                                                    'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                ]"
                                                            >
                                                                <SvgIcon
                                                                    class="w-5 h-5"
                                                                    aria-hidden="true"
                                                                    type="mdi"
                                                                    :path="mdiCheck"
                                                                />
                                                            </span>
                                                        </li>
                                                    </ComboboxOption>
                                                </ComboboxOptions>
                                            </div>
                                        </Combobox>
                                    </div>
                                    <div class="flex-1 form-control">
                                        <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.creator') }}
                                        </label>
                                        <Combobox as="div" v-model="search.character" class="mt-2" nullable>
                                            <div class="relative">
                                                <ComboboxButton as="div">
                                                    <ComboboxInput
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @change="queryChars = $event.target.value"
                                                        :display-value="(char: any) => `${char?.firstname} ${char?.lastname}`"
                                                        :placeholder="$t('common.creator')"
                                                    />
                                                </ComboboxButton>

                                                <ComboboxOptions
                                                    v-if="entriesChars.length > 0"
                                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm"
                                                >
                                                    <ComboboxOption
                                                        v-for="char in entriesChars"
                                                        :key="char.identifier"
                                                        :value="char"
                                                        as="char"
                                                        v-slot="{ active, selected }"
                                                    >
                                                        <li
                                                            :class="[
                                                                'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                active ? 'bg-primary-500' : '',
                                                            ]"
                                                        >
                                                            <span :class="['block truncate', selected && 'font-semibold']">
                                                                {{ char.firstname }} {{ char.lastname }}
                                                            </span>

                                                            <span
                                                                v-if="selected"
                                                                :class="[
                                                                    active ? 'text-neutral' : 'text-primary-500',
                                                                    'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                ]"
                                                            >
                                                                <SvgIcon
                                                                    class="w-5 h-5"
                                                                    aria-hidden="true"
                                                                    type="mdi"
                                                                    :path="mdiCheck"
                                                                />
                                                            </span>
                                                        </li>
                                                    </ComboboxOption>
                                                </ComboboxOptions>
                                            </div>
                                        </Combobox>
                                    </div>
                                    <div class="flex-1 form-control">
                                        <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.time_range') }}: {{ $t('common.from') }}
                                        </label>
                                        <div class="relative flex items-center mt-2">
                                            <input
                                                v-model="search.from"
                                                ref="searchInput"
                                                type="datetime-local"
                                                name="search"
                                                :placeholder="`${$t('common.time_range')} ${$t('common.from')}`"
                                                class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            />
                                        </div>
                                    </div>
                                    <div class="flex-1 form-control">
                                        <label for="search" class="block text-sm font-medium leading-6 text-neutral"
                                            >{{ $t('common.time_range') }}:
                                            {{ $t('common.to') }}
                                        </label>
                                        <div class="relative flex items-center mt-2">
                                            <input
                                                v-model="search.from"
                                                ref="searchInput"
                                                type="datetime-local"
                                                name="search"
                                                :placeholder="`${$t('common.time_range')} ${$t('common.to')}`"
                                                class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            />
                                        </div>
                                    </div>
                                </div>
                            </DisclosurePanel>
                        </Disclosure>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.document', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="documents && documents.length === 0"
                            :type="$t('common.document', 2)"
                            :focus="focusSearch"
                        />
                        <div v-else>
                            <ul class="flex flex-col">
                                <DocumentsListEntry v-for="doc in documents" :doc="doc" />
                            </ul>

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" class="mt-2" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

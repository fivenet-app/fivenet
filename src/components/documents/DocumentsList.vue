<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { BriefcaseIcon, CalendarIcon, DocumentMagnifyingGlassIcon, UserIcon } from '@heroicons/vue/20/solid';
import { CheckIcon } from '@heroicons/vue/24/solid';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/shared';
import { ref } from 'vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import TablePagination from '~/components/partials/TablePagination.vue';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { DocumentCategory } from '~~/gen/ts/resources/documents/category';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { ListDocumentsRequest } from '~~/gen/ts/services/docstore/docstore';
import TemplatesModal from './templates/TemplatesModal.vue';

const { $grpc } = useNuxtApp();

const search = ref<{ title: string; category?: DocumentCategory }>({ title: '' });
const pagination = ref<PaginationResponse>();
const offset = ref(0);

const entriesCategories = ref<DocumentCategory[]>([]);
const queryCategories = ref<string>('');

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

onMounted(async () => {
    findCategories();
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
                            <div class="flex-1 form-control">
                                <Combobox as="div" v-model="search.category" nullable>
                                    <div class="relative">
                                        <ComboboxButton as="div">
                                            <ComboboxInput
                                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @change="queryCategories = $event.target.value"
                                                :display-value="(category: any) => category?.name"
                                                placeholder="Category"
                                            />
                                        </ComboboxButton>

                                        <ComboboxOptions
                                            v-if="entriesCategories.length > 0"
                                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm"
                                        >
                                            <ComboboxOption
                                                v-for="category in entriesCategories"
                                                :key="category.id"
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
                                                        <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                    </span>
                                                </li>
                                            </ComboboxOption>
                                        </ComboboxOptions>
                                    </div>
                                </Combobox>
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
                        <button
                            v-else-if="documents && documents.length === 0"
                            type="button"
                            @click="focusSearch()"
                            class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                        >
                            <DocumentMagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                {{ $t('common.not_found', [$t('common.document', 2)]) }}
                                {{ $t('components.documents.document_list.no_documents_hint') }}
                            </span>
                        </button>
                        <div v-else>
                            <ul class="flex flex-col">
                                <li
                                    v-for="doc in documents"
                                    :key="doc.id"
                                    class="flex-initial my-1 rounded-lg hover:bg-base-800 bg-base-850"
                                >
                                    <NuxtLink
                                        :to="{
                                            name: 'documents-id',
                                            params: { id: doc.id },
                                        }"
                                    >
                                        <div class="mx-2 mt-1 mb-4">
                                            <div class="flex flex-row">
                                                <p class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
                                                    {{ doc.title }}
                                                </p>
                                                <p
                                                    class="inline-flex px-2 text-xs font-semibold leading-5 rounded-full bg-primary-100 text-primary-700 my-auto"
                                                    v-if="doc.state"
                                                >
                                                    {{ doc.state }}
                                                </p>
                                            </div>
                                            <div class="flex flex-row gap-2 text-base-200">
                                                <div class="flex flex-row items-center justify-start flex-1">
                                                    <UserIcon
                                                        class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                        aria-hidden="true"
                                                    />
                                                    {{ doc.creator?.firstname }},
                                                    {{ doc.creator?.lastname }}
                                                </div>
                                                <div class="flex flex-row items-center justify-center flex-1">
                                                    <BriefcaseIcon
                                                        class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                        aria-hidden="true"
                                                    />
                                                    {{ doc.creator?.jobLabel }}
                                                </div>
                                                <div class="flex flex-row items-center justify-end flex-1">
                                                    <CalendarIcon
                                                        class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                        aria-hidden="true"
                                                    />
                                                    <p>
                                                        {{ $t('common.created') }}
                                                        <time :datetime="$d(toDate(doc.createdAt)!, 'short')">
                                                            {{ useLocaleTimeAgo(toDate(doc.createdAt)!).value }}
                                                        </time>
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </NuxtLink>
                                </li>
                            </ul>

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" class="mt-2" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

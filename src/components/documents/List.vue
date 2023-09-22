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
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
} from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/shared';
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { useCompletorStore } from '~/store/completor';
import * as google_protobuf_timestamp_pb from '~~/gen/ts/google/protobuf/timestamp';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { Category } from '~~/gen/ts/resources/documents/category';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ListDocumentsRequest } from '~~/gen/ts/services/docstore/docstore';
import ListEntry from './ListEntry.vue';
import TemplatesModal from './templates/TemplatesModal.vue';

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const { t } = useI18n();

type OpenClose = { id: number; label: string; closed?: boolean };
const openclose: OpenClose[] = [
    { id: 0, label: t('common.not_selected') },
    { id: 1, label: t('common.open'), closed: false },
    { id: 2, label: t('common.close', 2), closed: true },
];

const query = ref<{
    title: string;
    category?: Category;
    character?: UserShort;
    from?: string;
    to?: string;
    closed?: OpenClose;
}>({
    title: '',
    closed: openclose[0],
});
const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const entriesCategories = ref<Category[]>([]);
const queryCategories = ref<string>('');
const entriesCitizens = ref<UserShort[]>([]);
const queryCitizens = ref<string>('');

const { data: documents, pending, refresh, error } = useLazyAsyncData(`documents-${offset.value}`, () => listDocuments());

async function listDocuments(): Promise<DocumentShort[]> {
    return new Promise(async (res, rej) => {
        const req: ListDocumentsRequest = {
            pagination: {
                offset: offset.value,
            },
            orderBy: [],
            search: query.value.title,
            categoryIds: [],
            creatorIds: [],
        };
        if (query.value.category) req.categoryIds.push(query.value.category.id);
        if (query.value.character) req.creatorIds.push(query.value.character.userId);
        if (query.value.from) {
            req.from = {
                timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(query.value.from)!),
            };
        }
        if (query.value.to) {
            req.to = {
                timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(query.value.to)!),
            };
        }
        if (query.value.closed) {
            if (query.value.closed !== undefined) {
                req.closed = query.value.closed.closed;
            }
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

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), { debounce: 600, maxWait: 1400 });

async function findCategories(): Promise<void> {
    entriesCategories.value = await completorStore.completeDocumentCategories(queryCategories.value);
}

watchDebounced(queryCategories, async () => findCategories(), {
    debounce: 600,
    maxWait: 1400,
});

watchDebounced(
    queryCitizens,
    async () =>
        (entriesCitizens.value = await completorStore.completeCitizens({
            search: queryCitizens.value,
        })),
    {
        debounce: 600,
        maxWait: 1400,
    },
);

onMounted(async () => findCategories());

const templatesOpen = ref(false);
</script>

<template>
    <TemplatesModal :open="templatesOpen" @close="templatesOpen = false" />

    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <label for="search" class="block mb-2 text-sm font-medium leading-6 text-neutral">
                            {{ $t('common.search') }}
                        </label>
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-1 form-control">
                                <input
                                    v-model="query.title"
                                    ref="searchInput"
                                    type="text"
                                    name="search"
                                    :placeholder="$t('common.title')"
                                    class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                />
                            </div>
                            <div class="flex-initial form-control" v-if="can('DocStoreService.CreateDocument')">
                                <button
                                    type="button"
                                    @click="templatesOpen = true"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                >
                                    {{ $t('common.create') }}
                                </button>
                            </div>
                            <div class="flex-initial" v-if="can('CompletorService.CompleteDocumentCategories')">
                                <NuxtLink
                                    :to="{ name: 'documents-categories' }"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                >
                                    {{ $t('common.category', 2) }}
                                </NuxtLink>
                            </div>
                            <div class="flex-initial" v-if="can('DocStoreService.ListTemplates')">
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
                                    <ChevronDownIcon
                                        :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                        aria-hidden="true"
                                    />
                                </span>
                            </DisclosureButton>
                            <DisclosurePanel class="mt-2 pr-4">
                                <div class="flex flex-row gap-2">
                                    <div class="flex-1 form-control">
                                        <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.category', 1) }}
                                        </label>
                                        <Combobox as="div" v-model="query.category" class="mt-2" nullable>
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
                                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
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
                                                                <CheckIcon class="w-5 h-5" aria-hidden="true" />
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
                                        <Combobox as="div" v-model="query.character" class="mt-2" nullable>
                                            <div class="relative">
                                                <ComboboxButton as="div">
                                                    <ComboboxInput
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @change="queryCitizens = $event.target.value"
                                                        :display-value="(char: any) => `${char?.firstname} ${char?.lastname}`"
                                                        :placeholder="$t('common.creator')"
                                                    />
                                                </ComboboxButton>

                                                <ComboboxOptions
                                                    v-if="entriesCitizens.length > 0"
                                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                >
                                                    <ComboboxOption
                                                        v-for="char in entriesCitizens"
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
                                                                <CheckIcon class="w-5 h-5" aria-hidden="true" />
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
                                                v-model="query.from"
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
                                                v-model="query.to"
                                                type="datetime-local"
                                                name="search"
                                                :placeholder="`${$t('common.time_range')} ${$t('common.to')}`"
                                                class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            />
                                        </div>
                                    </div>
                                    <div class="flex-1 form-control">
                                        <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.close', 2) }}?
                                        </label>
                                        <Listbox as="div" class="mt-2" v-model="query.closed">
                                            <div class="relative">
                                                <ListboxButton
                                                    class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                >
                                                    <span class="block truncate">
                                                        {{ query.closed?.label }}
                                                    </span>
                                                    <span
                                                        class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none"
                                                    >
                                                        <ChevronDownIcon class="w-5 h-5 text-gray-400" aria-hidden="true" />
                                                    </span>
                                                </ListboxButton>

                                                <transition
                                                    leave-active-class="transition duration-100 ease-in"
                                                    leave-from-class="opacity-100"
                                                    leave-to-class="opacity-0"
                                                >
                                                    <ListboxOptions
                                                        class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                    >
                                                        <ListboxOption
                                                            as="template"
                                                            v-for="st in openclose"
                                                            :key="st.id?.toString()"
                                                            :value="st"
                                                            v-slot="{ active, selected }"
                                                        >
                                                            <li
                                                                :class="[
                                                                    active ? 'bg-primary-500' : '',
                                                                    'text-neutral relative cursor-default select-none py-2 pl-8 pr-4',
                                                                ]"
                                                            >
                                                                <span
                                                                    :class="[
                                                                        selected ? 'font-semibold' : 'font-normal',
                                                                        'block truncate',
                                                                    ]"
                                                                    >{{ st.label }}</span
                                                                >

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
                                                        </ListboxOption>
                                                    </ListboxOptions>
                                                </transition>
                                            </div>
                                        </Listbox>
                                    </div>
                                </div>
                            </DisclosurePanel>
                        </Disclosure>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
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
                                <ListEntry v-for="doc in documents" :doc="doc" />
                            </ul>

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" class="mt-2" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

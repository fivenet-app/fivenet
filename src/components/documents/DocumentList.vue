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
import { watchDebounced } from '@vueuse/shared';
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { useCompletorStore } from '~/store/completor';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import { Category } from '~~/gen/ts/resources/documents/category';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ListDocumentsRequest, ListDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';
import DocumentListEntry from '~/components/documents/DocumentListEntry.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const { t } = useI18n();

type OpenClose = { id: number; label: string; closed?: boolean };
const openclose: OpenClose[] = [
    { id: 0, label: t('common.not_selected') },
    { id: 1, label: t('common.open', 2), closed: false },
    { id: 2, label: t('common.close', 2), closed: true },
];

const query = ref<{
    id?: string;
    title: string;
    character?: UserShort;
    from?: string;
    to?: string;
    closed?: boolean;
    category?: Category;
}>({
    title: '',
});
const offset = ref(0);

const queryClosed = ref<OpenClose>(openclose[0]);

const entriesCategories = ref<Category[]>([]);
const queryCategories = ref<string>('');

const entriesCitizens = ref<UserShort[]>([]);
const queryCitizens = ref<string>('');

const { data, pending, refresh, error } = useLazyAsyncData(`documents-${offset.value}`, () => listDocuments());

async function listDocuments(): Promise<ListDocumentsResponse> {
    const req: ListDocumentsRequest = {
        pagination: {
            offset: offset.value,
        },
        orderBy: [],
        search: query.value.title ?? '',
        categoryIds: [],
        creatorIds: [],
        documentIds: [],
    };
    if (query.value.category) {
        req.categoryIds.push(query.value.category.id);
    }
    if (query.value.character) {
        req.creatorIds.push(query.value.character.userId);
    }
    if (query.value.id) {
        const id = query.value.id.trim().replaceAll('-', '').replace(/\D/g, '');
        if (id.length > 0) {
            req.documentIds.push(id);
        }
    }
    if (query.value.from) {
        req.from = {
            timestamp: googleProtobufTimestamp.Timestamp.fromDate(fromString(query.value.from)!),
        };
    }
    if (query.value.to) {
        req.to = {
            timestamp: googleProtobufTimestamp.Timestamp.fromDate(fromString(query.value.to)!),
        };
    }
    if (query.value.closed !== undefined) {
        req.closed = query.value.closed;
    }

    try {
        const call = $grpc.getDocStoreClient().listDocuments(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
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

watch(queryClosed, () => (query.value.closed = queryClosed.value.closed));

onMounted(async () => findCategories());

const templatesOpen = ref(false);
</script>

<template>
    <TemplatesModal :open="templatesOpen" @close="templatesOpen = false" />

    <UDashboardToolbar>
        <template #default>
            <form class="w-full" @submit.prevent="refresh()">
                <div class="w-full flex flex-row gap-2">
                    <div class="flex-1">
                        <label for="search" class="mb-2 block text-sm font-medium leading-6 text-neutral">
                            {{ $t('common.search') }}
                        </label>
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-1">
                                <UInput
                                    ref="searchInput"
                                    v-model="query.title"
                                    type="text"
                                    name="search"
                                    :placeholder="$t('common.title')"
                                    class="hidden lg:block"
                                    block
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                    @keydown.esc="$event.target.blur()"
                                >
                                    <template #trailing>
                                        <UKbd value="/" />
                                    </template>
                                </UInput>
                            </div>
                        </div>
                    </div>
                </div>

                <Disclosure v-slot="{ open }" as="div" class="pt-2">
                    <DisclosureButton class="flex w-full items-start justify-between text-left text-sm text-neutral">
                        <span class="leading-7 text-accent-200">{{ $t('common.advanced_search') }}</span>
                        <span class="ml-6 flex h-7 items-center">
                            <ChevronDownIcon
                                :class="[open ? 'upsidedown' : '', 'size-5 transition-transform']"
                                aria-hidden="true"
                            />
                        </span>
                    </DisclosureButton>
                    <DisclosurePanel class="mt-2 pr-4">
                        <div class="flex flex-row flex-wrap gap-2">
                            <div class="flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.document') }} {{ $t('common.id') }}
                                </label>
                                <div class="relative mt-2">
                                    <input
                                        v-model="query.id"
                                        type="text"
                                        name="search"
                                        placeholder="DOC-..."
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.category', 1) }}
                                </label>
                                <Combobox v-model="query.category" as="div" class="mt-2" nullable>
                                    <div class="relative">
                                        <ComboboxButton as="div">
                                            <ComboboxInput
                                                autocomplete="off"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :display-value="(category: any) => category?.name"
                                                :placeholder="$t('common.category', 1)"
                                                @change="queryCategories = $event.target.value"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </ComboboxButton>

                                        <ComboboxOptions
                                            v-if="entriesCategories.length > 0"
                                            class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                        >
                                            <ComboboxOption
                                                v-for="category in entriesCategories"
                                                :key="category.id"
                                                v-slot="{ active, selected }"
                                                :value="category"
                                                as="category"
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
                                                        <CheckIcon class="size-5" aria-hidden="true" />
                                                    </span>
                                                </li>
                                            </ComboboxOption>
                                        </ComboboxOptions>
                                    </div>
                                </Combobox>
                            </div>
                            <div class="flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.creator') }}
                                </label>
                                <Combobox v-model="query.character" as="div" class="mt-2" nullable>
                                    <div class="relative">
                                        <ComboboxButton as="div">
                                            <ComboboxInput
                                                autocomplete="off"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :display-value="
                                                    (char: any) => `${char?.firstname} ${char?.lastname} (${char?.dateofbirth})`
                                                "
                                                :placeholder="$t('common.creator')"
                                                @change="queryCitizens = $event.target.value"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </ComboboxButton>

                                        <ComboboxOptions
                                            v-if="entriesCitizens.length > 0"
                                            class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                        >
                                            <ComboboxOption
                                                v-for="char in entriesCitizens"
                                                :key="char.identifier"
                                                v-slot="{ active, selected }"
                                                :value="char"
                                                as="char"
                                            >
                                                <li
                                                    :class="[
                                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                        active ? 'bg-primary-500' : '',
                                                    ]"
                                                >
                                                    <span :class="['block truncate', selected && 'font-semibold']">
                                                        {{ char.firstname }} {{ char.lastname }} ({{ char?.dateofbirth }})
                                                    </span>

                                                    <span
                                                        v-if="selected"
                                                        :class="[
                                                            active ? 'text-neutral' : 'text-primary-500',
                                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                        ]"
                                                    >
                                                        <CheckIcon class="size-5" aria-hidden="true" />
                                                    </span>
                                                </li>
                                            </ComboboxOption>
                                        </ComboboxOptions>
                                    </div>
                                </Combobox>
                            </div>
                        </div>
                        <div class="flex flex-row flex-wrap gap-2">
                            <div class="flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.close', 2) }}?
                                </label>
                                <Listbox v-model="queryClosed" as="div" class="mt-2">
                                    <div class="relative">
                                        <ListboxButton
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        >
                                            <span class="block truncate">
                                                {{ queryClosed?.label }}
                                            </span>
                                            <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                                <ChevronDownIcon class="size-5 text-gray-400" aria-hidden="true" />
                                            </span>
                                        </ListboxButton>

                                        <transition
                                            leave-active-class="transition duration-100 ease-in"
                                            leave-from-class="opacity-100"
                                            leave-to-class="opacity-0"
                                        >
                                            <ListboxOptions
                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                            >
                                                <ListboxOption
                                                    v-for="st in openclose"
                                                    v-slot="{ active, selected }"
                                                    :key="st.id"
                                                    as="template"
                                                    :value="st"
                                                >
                                                    <li
                                                        :class="[
                                                            active ? 'bg-primary-500' : '',
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
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
                                                            <CheckIcon class="size-5" aria-hidden="true" />
                                                        </span>
                                                    </li>
                                                </ListboxOption>
                                            </ListboxOptions>
                                        </transition>
                                    </div>
                                </Listbox>
                            </div>
                            <div class="flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.time_range') }}: {{ $t('common.from') }}
                                </label>
                                <div class="relative mt-2">
                                    <input
                                        v-model="query.from"
                                        type="datetime-local"
                                        name="search"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.from')}`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                            <div class="flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral"
                                    >{{ $t('common.time_range') }}:
                                    {{ $t('common.to') }}
                                </label>
                                <div class="relative mt-2">
                                    <input
                                        v-model="query.to"
                                        type="datetime-local"
                                        name="search"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.to')}`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                        </div>
                    </DisclosurePanel>
                </Disclosure>
            </form>
        </template>
    </UDashboardToolbar>

    <div class="mt-2 flow-root">
        <div class="-my-2 mx-0 overflow-x-auto">
            <div class="inline-block w-full max-w-full px-1 py-2 align-middle">
                <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.document', 2)])"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="data?.documents.length === 0"
                    :type="$t('common.document', 2)"
                    :focus="focusSearch"
                />
                <template v-else>
                    <ul role="list" class="flex flex-col">
                        <DocumentListEntry v-for="doc in data?.documents" :key="doc.id" :doc="doc" />
                    </ul>

                    <TablePagination
                        class="mt-2"
                        :pagination="data?.pagination"
                        :refresh="refresh"
                        @offset-change="offset = $event"
                    />
                </template>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import {
    Combobox,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Dialog,
    DialogPanel,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';

import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { onClickOutside, useMagicKeys, watchDebounced, whenever } from '@vueuse/core';
import {
    AccountIcon,
    AccountMultipleIcon,
    BriefcaseIcon,
    CarIcon,
    CogIcon,
    FileDocumentIcon,
    FileDocumentMultipleIcon,
    GlobeModelIcon,
    HomeIcon,
    MagnifyIcon,
    MapIcon,
    RefreshIcon,
} from 'mdi-vue3';
import { DefineComponent } from 'vue';
import '~/assets/css/command-palette.scss';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ListDocumentsRequest } from '~~/gen/ts/services/docstore/docstore';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const open = ref(false);
const target = ref(null);

const loading = ref(false);
const rawQuery = ref('');
const query = computed(() => rawQuery.value.replace(/^[@#]/, ''));

const keys = useMagicKeys({
    passive: false,
    onEventFired(e) {
        if (
            ((e.metaKey || e.ctrlKey) && e.key === 'k' && e.type === 'keydown') ||
            (e.shiftKey && e.key === '/' && e.type === 'keydown')
        ) {
            e.preventDefault();
            open.value = true;
        }
    },
});
const Escape = keys['Escape'];

whenever(Escape, () => {
    open.value = false;
});

onClickOutside(target, () => {
    open.value = false;
});

type Item = {
    id: number;
    name: string;
    prefix?: string;
    icon?: DefineComponent;
    category: string;
    permission?: string;
    action: () => any;
};

const items = [
    {
        id: 1,
        name: t('commandpalette.groups.shortcuts.goto', [t('common.citizen', 1), 'CIT-...']),
        icon: markRaw(AccountIcon),
        category: 'shortcuts',
        permission: 'CitizenStoreService.GetCitizen',
        prefix: 'CIT-',
        action: () => {
            const id = query.value.substring(query.value.indexOf('-') + 1);
            if (id.length > 0) {
                navigateTo({ name: 'citizens-id', params: { id: id } });
                open.value = false;
            }
        },
    },
    {
        id: 2,
        name: t('commandpalette.groups.shortcuts.goto', [t('common.document', 1), 'DOC-...']),
        icon: markRaw(FileDocumentIcon),
        category: 'shortcuts',
        permission: 'DocStoreService.GetDocument',
        prefix: 'DOC-',
        action: () => {
            const id = query.value.substring(query.value.indexOf('-') + 1);
            if (id.length > 0) {
                navigateTo({ name: 'documents-id', params: { id: id } });
                open.value = false;
            }
        },
    },
    // Pages
    {
        id: 9,
        name: t('common.overview'),
        href: { name: 'overview' },
        icon: markRaw(HomeIcon),
        category: 'pages',
    },
    {
        id: 10,
        name: t('common.citizen'),
        action: () => {
            navigateTo({ name: 'citizens' });
            open.value = false;
        },
        permission: 'CitizenStoreService.ListCitizens',
        icon: markRaw(AccountMultipleIcon),
        category: 'pages',
    },
    {
        id: 11,
        name: t('common.vehicle'),
        action: () => {
            navigateTo({ name: 'vehicles' });
            open.value = false;
        },
        permission: 'DMVService.ListVehicles',
        icon: markRaw(CarIcon),
        category: 'pages',
    },
    {
        id: 12,
        name: t('common.document'),
        action: () => {
            navigateTo({ name: 'documents' });
            open.value = false;
        },
        permission: 'DocStoreService.ListDocuments',
        icon: markRaw(FileDocumentMultipleIcon),
        category: 'pages',
    },
    {
        id: 13,
        name: t('common.job'),
        action: () => {
            navigateTo({ name: 'jobs' });
            open.value = false;
        },
        permission: 'Jobs.View',
        icon: markRaw(BriefcaseIcon),
        category: 'pages',
    },
    {
        id: 14,
        name: t('common.livemap'),
        action: () => {
            navigateTo({ name: 'livemap' });
            open.value = false;
        },
        permission: 'LivemapperService.Stream',
        icon: markRaw(MapIcon),
        category: 'pages',
    },
    {
        id: 15,
        name: t('common.control_panel'),
        action: () => {
            navigateTo({ name: 'rector' });
            open.value = false;
        },
        permission: 'RectorService.GetRoles',
        icon: markRaw(CogIcon),
        category: 'pages',
    },
] as Item[];

const filteredItems = computed<Item[]>(() =>
    query.value === ''
        ? items
        : items.filter((item) => {
              if (item.prefix) {
                  if (query.value.toLowerCase().startsWith(item.prefix.toLowerCase())) {
                      return true;
                  }
              }
              return item.name.toLowerCase().includes(query.value.toLowerCase());
          }),
);

const citizens = ref<UserShort[]>([]);

const filteredCitizens = computed(() =>
    citizens.value.filter((citizen) =>
        `${citizen.firstname} ${citizen.lastname}`.toLowerCase().includes(query.value.toLowerCase()),
    ),
);

async function listCitizens(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const req = {
                pagination: {
                    offset: 0n,
                    pageSize: 7n,
                },
                searchName: query.value,
            };

            const call = $grpc.getCitizenStoreClient().listCitizens(req);
            const { response } = await call;

            citizens.value = response.users;
            loading.value = false;
            return res();
        } catch (e) {
            loading.value = false;
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const documents = ref<DocumentShort[]>([]);

const filteredDocuments = computed(() =>
    documents.value.filter((document) => document.title.toLowerCase().includes(query.value.toLowerCase())),
);

async function listDocuments(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req: ListDocumentsRequest = {
            pagination: {
                offset: 0n,
                pageSize: 7n,
            },
            orderBy: [],
            search: query.value,
            categoryIds: [],
            creatorIds: [],
        };

        try {
            const call = $grpc.getDocStoreClient().listDocuments(req);
            const { response } = await call;

            documents.value = response.documents;
            loading.value = false;
            return res();
        } catch (e) {
            loading.value = false;
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(query, async () => (loading.value = true));
watchDebounced(
    query,
    async () => {
        if (query.value.length > 1) {
            if (rawQuery.value.startsWith('@')) {
                listCitizens();
            } else if (rawQuery.value.startsWith('#')) {
                listDocuments();
            }
        } else {
            citizens.value.length = 0;
            documents.value.length = 0;
            loading.value = false;
        }
    },
    { debounce: 500, maxWait: 1500 },
);

type Groups = { [key: string]: Item[] };
const groups = computed<Groups>(() =>
    filteredItems.value.reduce((groups, item) => {
        return { ...groups, [item.category]: [...((groups as Groups)[item.category] || []), item] };
    }, {}),
);

async function onSelect(item: any): Promise<any> {
    if ('userId' in item) {
        navigateTo({ name: 'citizens-id', params: { id: item.userId } });
        open.value = false;
        return;
    } else if ('title' in item && 'id' in item) {
        navigateTo({ name: 'documents-id', params: { id: item.id } });
        open.value = false;
        return;
    } else if ('action' in item) {
        return item.action();
    }
}
</script>

<template>
    <TransitionRoot :show="open" as="template" @after-leave="rawQuery = ''" appear>
        <Dialog as="div" class="relative z-10" @close="open = false">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-gray-900 bg-opacity-25 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto p-4 sm:p-6 md:p-20">
                <TransitionChild
                    as="template"
                    enter="ease-out duration-300"
                    enter-from="opacity-0 scale-95"
                    enter-to="opacity-100 scale-100"
                    leave="ease-in duration-200"
                    leave-from="opacity-100 scale-100"
                    leave-to="opacity-0 scale-95"
                >
                    <DialogPanel
                        class="mx-auto max-w-3xl transform overflow-hidden rounded-xl bg-gray-900 shadow-2xl ring-1 ring-black ring-opacity-5 transition-all"
                    >
                        <Combobox @update:modelValue="onSelect">
                            <div class="relative">
                                <MagnifyIcon
                                    class="pointer-events-none absolute left-4 top-3.5 h-5 w-5 text-gray-400"
                                    aria-hidden="true"
                                />
                                <ComboboxInput
                                    class="h-12 w-full border-0 bg-transparent pl-11 pr-4 text-gray-200 placeholder:text-gray-400 focus:ring-0 sm:text-sm"
                                    placeholder="Search..."
                                    @change="rawQuery = $event.target.value"
                                />
                            </div>

                            <div
                                v-if="rawQuery === ''"
                                class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14"
                            >
                                <GlobeModelIcon class="mx-auto h-6 w-6 text-gray-400" aria-hidden="true" />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('commandpalette.input.title') }}
                                </p>
                                <p class="mt-2 text-gray-500">
                                    {{ $t('commandpalette.input.content') }}
                                </p>
                            </div>

                            <ComboboxOptions static class="max-h-80 scroll-pb-2 scroll-pt-11 space-y-2 overflow-y-auto pb-2">
                                <li v-if="filteredCitizens.length > 0" class="p-2">
                                    <h2 class="text-xs font-semibold text-gray-900">{{ $t('common.citizen', 2) }}</h2>
                                    <ul class="-mx-2 mt-2 text-sm text-gray-400">
                                        <ComboboxOption
                                            v-for="citizen in filteredCitizens"
                                            :key="citizen.userId"
                                            :value="citizen"
                                            as="template"
                                            v-slot="{ active }"
                                        >
                                            <li
                                                :class="[
                                                    'flex cursor-default select-none items-center px-4 py-2',
                                                    active && 'bg-primary-600 text-white',
                                                ]"
                                            >
                                                <span class="ml-3 flex-auto truncate">
                                                    {{ citizen.firstname }} {{ citizen.lastname }}
                                                </span>
                                            </li>
                                        </ComboboxOption>
                                    </ul>
                                </li>
                                <li v-else-if="filteredDocuments.length > 0" class="p-2">
                                    <h2 class="text-xs font-semibold text-gray-900">{{ $t('common.document', 2) }}</h2>
                                    <ul class="-mx-2 mt-2 text-sm text-gray-400">
                                        <ComboboxOption
                                            v-for="document in filteredDocuments"
                                            :key="document.id.toString()"
                                            :value="document"
                                            as="template"
                                            v-slot="{ active }"
                                        >
                                            <li
                                                :class="[
                                                    'flex cursor-default select-none items-center px-4 py-2',
                                                    active && 'bg-gray-800 text-white',
                                                ]"
                                            >
                                                <span class="ml-3 flex-auto truncate">{{ document.title }}</span>
                                            </li>
                                        </ComboboxOption>
                                    </ul>
                                </li>
                                <li v-else v-for="[category, items] in Object.entries(groups)" :key="category" class="p-2">
                                    <h2 class="px-4 py-2.5 text-xs font-semibold text-gray-200">
                                        {{ $t(`commandpalette.groups.${category}.label`) }}
                                    </h2>
                                    <ul class="-mx-2 mt-2 text-sm text-gray-400">
                                        <ComboboxOption
                                            v-for="item in items"
                                            :key="item.id"
                                            :value="item"
                                            as="template"
                                            v-slot="{ active }"
                                            v-can="item.permission"
                                        >
                                            <li
                                                :class="[
                                                    'flex cursor-default select-none items-center rounded-md px-3 py-2',
                                                    active && 'bg-gray-800 text-white',
                                                ]"
                                            >
                                                <component
                                                    v-if="item.icon"
                                                    :is="item.icon"
                                                    :class="['h-6 w-6 flex-none', active ? 'text-white' : 'text-gray-500']"
                                                    aria-hidden="true"
                                                />
                                                <span class="ml-3 flex-auto truncate">{{ item.name }}</span>
                                                <span v-if="active" class="ml-3 flex-none text-gray-200">Jump to...</span>
                                            </li>
                                        </ComboboxOption>
                                    </ul>
                                </li>
                            </ComboboxOptions>

                            <div
                                v-if="(rawQuery.startsWith('@') || rawQuery.startsWith('#')) && loading"
                                class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14"
                            >
                                <RefreshIcon class="mx-auto h-6 w-6 text-gray-400 animate-spin" aria-hidden="true" />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('common.loading', [$t('common.result', 2)]) }}
                                </p>
                            </div>
                            <div
                                v-else-if="
                                    rawQuery !== '' &&
                                    filteredItems.length === 0 &&
                                    filteredCitizens.length === 0 &&
                                    filteredDocuments.length === 0
                                "
                                class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14"
                            >
                                <GlobeModelIcon class="mx-auto h-6 w-6 text-gray-400" aria-hidden="true" />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('commandpalette.empty.title') }}
                                </p>
                                <p class="mt-2 text-gray-500">
                                    {{ $t('commandpalette.empty.content') }}
                                </p>
                            </div>
                            <div class="flex flex-wrap items-center bg-gray-50 px-4 py-2.5 text-xs text-gray-700">
                                Type
                                <kbd
                                    :class="[
                                        'mx-1 flex h-5 w-5 items-center justify-center rounded border bg-white font-semibold sm:mx-2',
                                        rawQuery.startsWith('@')
                                            ? 'border-primary-600 text-primary-600'
                                            : 'border-gray-400 text-gray-900',
                                    ]"
                                    >@</kbd
                                >
                                <span class="sm:hidden">for citizens,</span>
                                <span class="hidden sm:inline">to search citizens,</span>
                                <kbd
                                    :class="[
                                        'mx-1 flex h-5 w-5 items-center justify-center rounded border bg-white font-semibold sm:mx-2',
                                        rawQuery.startsWith('#')
                                            ? 'border-primary-600 text-primary-600'
                                            : 'border-gray-400 text-gray-900',
                                    ]"
                                    >#</kbd
                                >
                                for documents.
                                <!-- <kbd :class="['mx-1 flex h-5 w-5 items-center justify-center rounded border bg-white font-semibold sm:mx-2', query === '?' ? 'border-primary-600 text-primary-600' : 'border-gray-400 text-gray-900']">?</kbd>
                                for help.-->
                            </div>
                        </Combobox>
                    </DialogPanel>
                </TransitionChild>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

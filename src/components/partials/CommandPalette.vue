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
import { RpcError } from '@protobuf-ts/runtime-rpc';
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
import { type DefineComponent } from 'vue';
import '~/assets/css/command-palette.scss';
import { toggleTablet } from '~/composables/nui';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { UserShort } from '~~/gen/ts/resources/users/users';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import type { Perms } from '~~/gen/ts/perms';

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
        if ((e.metaKey || e.ctrlKey) && e.key === 'k' && e.type === 'keydown') {
            e.preventDefault();
            open.value = true;
        }
    },
});
/* eslint dot-notation: "off" */
const Escape = keys['Escape'];

whenever(Escape, async () => {
    if (!open.value) {
        toggleTablet(false);
    }
    open.value = false;
});

onClickOutside(target, async () => {
    open.value = false;
});

type Item = {
    id: number;
    name: string;
    prefix?: string;
    icon?: DefineComponent;
    category: string;
    permission?: Perms;
    action: () => any;
};

const items = [
    {
        id: -1,
        name: t('commandpalette.groups.shortcuts.goto', { key: t('common.citizen', 1), prefix: 'CIT-...' }),
        icon: markRaw(AccountIcon),
        category: 'shortcuts',
        permission: 'CitizenStoreService.GetUser',
        prefix: 'CIT-',
        action: () => {
            const id = query.value.substring(query.value.indexOf('-') + 1);
            if (id.length > 0) {
                navigateTo({ name: 'citizens-id', params: { id } });
                open.value = false;
            }
        },
    },
    {
        id: -2,
        name: t('commandpalette.groups.shortcuts.goto', { key: t('common.document', 1), prefix: 'DOC-...' }),
        icon: markRaw(FileDocumentIcon),
        category: 'shortcuts',
        permission: 'DocStoreService.GetDocument',
        prefix: 'DOC-',
        action: () => {
            const id = query.value.substring(query.value.indexOf('-') + 1);
            if (id.length > 0) {
                navigateTo({ name: 'documents-id', params: { id } });
                open.value = false;
            }
        },
    },
    // Pages
    {
        id: -9,
        name: t('common.overview'),
        action: () => {
            navigateTo({ name: 'overview' });
            open.value = false;
        },
        icon: markRaw(HomeIcon),
        category: 'pages',
    },
    {
        id: -10,
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
        id: -11,
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
        id: -12,
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
        id: -13,
        name: t('common.job'),
        action: () => {
            navigateTo({ name: 'jobs-overview' });
            open.value = false;
        },
        permission: 'JobsService.ListColleagues',
        icon: markRaw(BriefcaseIcon),
        category: 'pages',
    },
    {
        id: -14,
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
        id: -15,
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

async function listCitizens(): Promise<void> {
    try {
        const call = $grpc.getCitizenStoreClient().listCitizens({
            pagination: {
                offset: 0n,
                pageSize: 7n,
            },
            searchName: query.value,
        });
        const { response } = await call;

        citizens.value = response.users;
        loading.value = false;
    } catch (e) {
        loading.value = false;
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const documents = ref<DocumentShort[]>([]);

async function listDocuments(): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().listDocuments({
            pagination: {
                offset: 0n,
                pageSize: 7n,
            },
            orderBy: [],
            search: query.value,
            categoryIds: [],
            creatorIds: [],
            documentIds: [],
        });
        const { response } = await call;

        documents.value = response.documents;
        loading.value = false;
    } catch (e) {
        loading.value = false;
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function resetLists(): void {
    citizens.value.length = 0;
    documents.value.length = 0;
    loading.value = false;
}

watch(query, async () => (loading.value = true));
watchDebounced(
    rawQuery,
    async () => {
        if (query.value.length > 1) {
            if (rawQuery.value.startsWith('@')) {
                if (documents.value.length > 0) documents.value.length = 0;
                listCitizens();
            } else if (rawQuery.value.startsWith('#')) {
                if (citizens.value.length > 0) citizens.value.length = 0;
                listDocuments();
            }
        } else {
            resetLists();
        }
    },
    { debounce: 400, maxWait: 1250 },
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
    } else if ('title' in item && 'id' in item) {
        navigateTo({ name: 'documents-id', params: { id: item.id } });
        open.value = false;
    } else if ('action' in item) {
        return item.action();
    }
}
</script>

<template>
    <TransitionRoot :show="open" as="template" appear @after-leave="rawQuery = ''">
        <Dialog as="div" class="relative z-30" @close="open = false">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-primary-900 bg-opacity-25 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto p-4 sm:p-6 md:p-20">
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
                        class="mx-auto max-w-3xl transform overflow-hidden rounded-xl bg-primary-900 shadow-2xl ring-1 ring-black ring-opacity-5 transition-all"
                    >
                        <Combobox @update:model-value="onSelect">
                            <div class="relative">
                                <MagnifyIcon
                                    class="pointer-events-none absolute left-4 top-3.5 h-5 w-5 text-gray-400"
                                    aria-hidden="true"
                                />
                                <ComboboxInput
                                    autocomplete="off"
                                    class="h-12 w-full border-0 bg-transparent pl-11 pr-4 text-gray-200 placeholder:text-gray-400 focus:ring-0 sm:text-sm"
                                    placeholder="Search..."
                                    @change="rawQuery = $event.target.value"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </div>

                            <div
                                v-if="rawQuery === ''"
                                class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14"
                            >
                                <GlobeModelIcon class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('commandpalette.input.title') }}
                                </p>
                                <p class="mt-2 text-gray-500">
                                    {{ $t('commandpalette.input.content') }}
                                </p>
                            </div>

                            <ComboboxOptions static class="max-h-80 scroll-pb-2 scroll-pt-11 space-y-2 overflow-y-auto pb-2">
                                <li v-if="citizens.length > 0" class="p-2">
                                    <h2 class="text-xs font-semibold text-gray-900">{{ $t('common.citizen', 2) }}</h2>
                                    <ul class="-mx-2 mt-2 text-sm text-gray-400">
                                        <ComboboxOption
                                            v-for="citizen in citizens"
                                            :key="citizen.userId"
                                            v-slot="{ active }"
                                            :value="citizen"
                                            as="template"
                                        >
                                            <li
                                                :class="[
                                                    'flex cursor-default select-none items-center px-4 py-2',
                                                    active && 'bg-primary-500 text-neutral',
                                                ]"
                                            >
                                                <span class="ml-3 flex-auto truncate">
                                                    {{ citizen.firstname }} {{ citizen.lastname }}
                                                </span>
                                                <IDCopyBadge :id="citizen.userId" prefix="CIT" class="self-end" />
                                            </li>
                                        </ComboboxOption>
                                    </ul>
                                </li>
                                <li v-else-if="documents.length > 0" class="p-2">
                                    <h2 class="text-xs font-semibold text-gray-900">{{ $t('common.document', 2) }}</h2>
                                    <ul class="-mx-2 mt-2 text-sm text-gray-400">
                                        <ComboboxOption
                                            v-for="document in documents"
                                            :key="document.id"
                                            v-slot="{ active }"
                                            :value="document"
                                            as="template"
                                        >
                                            <li
                                                :class="[
                                                    'flex cursor-default select-none items-center px-4 py-2',
                                                    active && 'bg-gray-800 text-neutral',
                                                ]"
                                            >
                                                <span class="ml-3 flex-auto truncate">{{ document.title }}</span>
                                                <IDCopyBadge :id="document.id" prefix="DOC" class="self-end" />
                                            </li>
                                        </ComboboxOption>
                                    </ul>
                                </li>
                                <li v-else class="p-2">
                                    <template v-for="[category, cItems] in Object.entries(groups)" :key="category">
                                        <h2 class="px-4 py-2.5 text-xs font-semibold text-gray-200">
                                            {{ $t(`commandpalette.groups.${category}.label`) }}
                                        </h2>
                                        <ul class="-mx-2 mt-2 text-sm text-gray-400">
                                            <ComboboxOption
                                                v-for="item in cItems.filter(
                                                    (e) => e.permission === undefined || can(e.permission),
                                                )"
                                                :key="item.id"
                                                v-slot="{ active }"
                                                :value="item"
                                                as="template"
                                            >
                                                <li
                                                    :class="[
                                                        'flex cursor-default select-none items-center rounded-md px-3 py-2',
                                                        active && 'bg-gray-800 text-neutral',
                                                    ]"
                                                >
                                                    <component
                                                        :is="item.icon"
                                                        v-if="item.icon"
                                                        :class="[
                                                            'h-5 w-5 flex-none',
                                                            active ? 'text-neutral' : 'text-gray-500',
                                                        ]"
                                                        aria-hidden="true"
                                                    />
                                                    <span class="ml-3 flex-auto truncate">{{ item.name }}</span>
                                                    <span v-if="active" class="ml-3 flex-none text-gray-200">Jump to...</span>
                                                </li>
                                            </ComboboxOption>
                                        </ul>
                                    </template>
                                </li>
                            </ComboboxOptions>

                            <div
                                v-if="(rawQuery.startsWith('@') || rawQuery.startsWith('#')) && loading"
                                class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14"
                            >
                                <RefreshIcon class="mx-auto h-5 w-5 animate-spin text-gray-400" aria-hidden="true" />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('common.loading', [$t('common.result', 2)]) }}
                                </p>
                            </div>
                            <div
                                v-else-if="
                                    rawQuery !== '' &&
                                    filteredItems.length === 0 &&
                                    citizens.length === 0 &&
                                    documents.length === 0
                                "
                                class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14"
                            >
                                <GlobeModelIcon class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('commandpalette.empty.title') }}
                                </p>
                                <p class="mt-2 text-gray-500">
                                    {{ $t('commandpalette.empty.content') }}
                                </p>
                            </div>
                            <I18nT
                                keypath="commandpalette.footer"
                                tag="div"
                                class="flex flex-wrap items-center bg-gray-50 px-4 py-2.5 text-xs text-gray-700"
                            >
                                <template #key1>
                                    <kbd
                                        :class="[
                                            'mx-1 flex h-5 w-5 items-center justify-center rounded border bg-neutral font-semibold sm:mx-2',
                                            rawQuery.startsWith('@')
                                                ? 'border-primary-600 text-primary-600'
                                                : 'border-gray-400 text-gray-900',
                                        ]"
                                        >@</kbd
                                    >
                                </template>
                                <template #key2>
                                    <kbd
                                        :class="[
                                            'mx-1 flex h-5 w-5 items-center justify-center rounded border bg-neutral font-semibold sm:mx-2',
                                            rawQuery.startsWith('#')
                                                ? 'border-primary-600 text-primary-600'
                                                : 'border-gray-400 text-gray-900',
                                        ]"
                                        >#</kbd
                                    >
                                </template>
                            </I18nT>
                        </Combobox>
                    </DialogPanel>
                </TransitionChild>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

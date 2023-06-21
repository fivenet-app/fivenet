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
import SvgIcon from '@jamescoyle/vue-icon';
import {
    mdiAccount,
    mdiAccountMultiple,
    mdiBriefcase,
    mdiCar,
    mdiCog,
    mdiFileDocument,
    mdiFileDocumentMultiple,
    mdiGlobeModel,
    mdiHome,
    mdiMagnify,
    mdiMap,
} from '@mdi/js';
import { onClickOutside, useMagicKeys, whenever } from '@vueuse/core';
import '~/assets/css/command-palette.scss';

const { t } = useI18n();

const open = ref(false);
const query = ref('');
const target = ref(null);

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
    icon?: string;
    category: string;
    permission?: string;
    action: () => any;
};

const items = [
    {
        id: 1,
        name: t('commandpalette.groups.shortcuts.goto', [t('common.citizen', 1), 'CIT-...']),
        icon: mdiAccount,
        category: 'shortcuts',
        permission: 'CitizenStoreService.GetCitizen',
        prefix: 'CIT-',
        action: () => {
            const id = query.value.substring(query.value.indexOf('-') + 1);
            if (id.length > 0) {
                navigateTo({ name: 'citizens-id', params: { id: id } });
                open.value = false;
            } else {
            }
        },
    },
    {
        id: 2,
        name: t('commandpalette.groups.shortcuts.goto', [t('common.document', 1), 'DOC-...']),
        icon: mdiFileDocument,
        category: 'shortcuts',
        permission: 'DocStoreService.GetDocument',
        prefix: 'DOC-',
        action: () => {
            const id = query.value.substring(query.value.indexOf('-') + 1);
            if (id.length > 0) {
                navigateTo({ name: 'documents-id', params: { id: id } });
                open.value = false;
            } else {
            }
        },
    },
    // Pages
    {
        id: 4,
        name: t('common.overview'),
        href: { name: 'overview' },
        icon: mdiHome,
        category: 'pages',
    },
    {
        id: 5,
        name: t('common.citizen'),
        action: () => {
            navigateTo({ name: 'citizens' });
            open.value = false;
        },
        permission: 'CitizenStoreService.ListCitizens',
        icon: mdiAccountMultiple,
        category: 'pages',
    },
    {
        id: 6,
        name: t('common.vehicle'),
        action: () => {
            navigateTo({ name: 'vehicles' });
            open.value = false;
        },
        permission: 'DMVService.ListVehicles',
        icon: mdiCar,
        category: 'pages',
    },
    {
        id: 7,
        name: t('common.document'),
        action: () => {
            navigateTo({ name: 'documents' });
            open.value = false;
        },
        permission: 'DocStoreService.ListDocuments',
        icon: mdiFileDocumentMultiple,
        category: 'pages',
    },
    {
        id: 8,
        name: t('common.job'),
        action: () => {
            navigateTo({ name: 'jobs' });
            open.value = false;
        },
        permission: 'Jobs.View',
        icon: mdiBriefcase,
        category: 'pages',
    },
    {
        id: 9,
        name: t('common.livemap'),
        action: () => {
            navigateTo({ name: 'livemap' });
            open.value = false;
        },
        permission: 'LivemapperService.Stream',
        icon: mdiMap,
        category: 'pages',
    },
    {
        id: 10,
        name: t('common.control_panel'),
        action: () => {
            navigateTo({ name: 'rector' });
            open.value = false;
        },
        permission: 'RectorService.GetRoles',
        icon: mdiCog,
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
          })
);

type Groups = { [key: string]: Item[] };
const groups = computed<Groups>(() =>
    filteredItems.value.reduce((groups, item) => {
        return { ...groups, [item.category]: [...((groups as Groups)[item.category] || []), item] };
    }, {})
);

async function onSelect(item: Item): Promise<any> {
    return item.action();
}
</script>

<template>
    <TransitionRoot :show="open" as="template" @after-leave="query = ''" appear>
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
                                <SvgIcon
                                    type="mdi"
                                    :path="mdiMagnify"
                                    class="pointer-events-none absolute left-4 top-3.5 h-5 w-5 text-gray-400"
                                    aria-hidden="true"
                                />
                                <ComboboxInput
                                    class="h-12 w-full border-0 bg-transparent pl-11 pr-4 text-gray-200 placeholder:text-gray-400 focus:ring-0 sm:text-sm"
                                    placeholder="Search..."
                                    @change="query = $event.target.value"
                                />
                            </div>

                            <div v-if="query === ''" class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14">
                                <SvgIcon
                                    type="mdi"
                                    :path="mdiGlobeModel"
                                    class="mx-auto h-6 w-6 text-gray-400"
                                    aria-hidden="true"
                                />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('commandpalette.input.title') }}
                                </p>
                                <p class="mt-2 text-gray-500">
                                    {{ $t('commandpalette.input.content') }}
                                </p>
                            </div>

                            <ComboboxOptions
                                v-if="filteredItems.length > 0"
                                static
                                class="max-h-80 scroll-pb-2 scroll-pt-11 space-y-2 overflow-y-auto pb-2"
                            >
                                <li v-for="[category, items] in Object.entries(groups)" :key="category">
                                    <h2 class="px-4 py-2.5 text-xs font-semibold text-gray-200">
                                        {{ $t(`commandpalette.groups.${category}.label`) }}
                                    </h2>
                                    <ul class="mt-2 text-sm text-gray-400">
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
                                                <SvgIcon
                                                    v-if="item.icon"
                                                    type="mdi"
                                                    :path="item.icon"
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
                                v-if="query !== '' && filteredItems.length === 0"
                                class="border-t border-gray-100 px-6 py-14 text-center text-sm sm:px-14"
                            >
                                <SvgIcon
                                    type="mdi"
                                    :path="mdiGlobeModel"
                                    class="mx-auto h-6 w-6 text-gray-400"
                                    aria-hidden="true"
                                />
                                <p class="mt-4 font-semibold text-gray-200">
                                    {{ $t('commandpalette.empty.title') }}
                                </p>
                                <p class="mt-2 text-gray-500">
                                    {{ $t('commandpalette.empty.content') }}
                                </p>
                            </div>
                            <!-- <div class="flex flex-wrap items-center bg-gray-50 px-4 py-2.5 text-xs text-gray-700">
                                Type
                                <kbd
                                    :class="[
                                        'mx-1 flex h-5 w-5 items-center justify-center rounded border bg-gray-900 font-semibold sm:mx-2',
                                        query === '?' ? 'border-primary-600 text-primary-600' : 'border-gray-400 text-gray-200',
                                    ]"
                                    >?</kbd
                                >
                                for help.
                            </div> -->
                        </Combobox>
                    </DialogPanel>
                </TransitionChild>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

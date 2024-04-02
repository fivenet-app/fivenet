<script lang="ts" setup>
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import type { RoutePathSchema } from '@typed-router';
import { useAuthStore } from '~/store/auth';
import { JOB_THEME_KEY, availableThemes, useSettingsStore } from '~/store/settings';
import type { Perms } from '~~/gen/ts/perms';
import GenericContainerPanel from '~/components/partials/elements/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/elements/GenericContainerPanelEntry.vue';

const { t } = useI18n();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const settings = useSettingsStore();
const { startpage, documents, streamerMode, audio } = storeToRefs(settings);

const homepages: { name: string; path: RoutePathSchema; permission?: Perms }[] = [
    { name: 'common.home', path: '/overview' },
    { name: 'pages.citizens.title', path: '/citizens', permission: 'CitizenStoreService.ListCitizens' },
    { name: 'pages.vehicles.title', path: '/vehicles', permission: 'DMVService.ListVehicles' },
    { name: 'pages.documents.title', path: '/documents', permission: 'DocStoreService.ListDocuments' },
    { name: 'pages.jobs.overview.title', path: '/jobs/overview', permission: 'JobsService.ListColleagues' },
    { name: 'common.livemap', path: '/livemap', permission: 'LivemapperService.Stream' },
    { name: 'common.dispatch_center', path: '/centrum', permission: 'CentrumService.TakeControl' },
];

const selectedHomepage = ref<(typeof homepages)[0]>();
watch(selectedHomepage, () => (startpage.value = selectedHomepage.value?.path ?? '/overview'));

onBeforeMount(async () => {
    selectedHomepage.value = homepages.find((h) => h.path === startpage.value);
});

const darkModeActive = ref(documents.value.editorTheme === 'dark');

watch(darkModeActive, async () => {
    if (darkModeActive.value) {
        documents.value.editorTheme = 'dark';
    } else {
        documents.value.editorTheme = 'default';
    }
});

const themes = [
    {
        name: t('components.auth.settings_panel.app_theme.job_default_theme'),
        key: JOB_THEME_KEY,
    },
    ...availableThemes,
];
</script>

<template>
    <GenericContainerPanel>
        <template #title>
            {{ $t('components.auth.settings_panel.title') }}
        </template>
        <template #description>
            {{ $t('components.auth.settings_panel.subtitle') }}
        </template>
        <template #default>
            <GenericContainerPanelEntry>
                <template #title>
                    {{ $t('common.theme') }}
                </template>
                <template #default>
                    <Listbox v-model="settings.theme" as="div">
                        <div class="relative">
                            <ListboxButton
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            >
                                <span class="block truncate">
                                    {{ themes.find((t) => t.key === settings.theme)?.name }}
                                </span>
                                <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                    <ChevronDownIcon class="size-5 text-gray-400" />
                                </span>
                            </ListboxButton>

                            <transition
                                leave-active-class="transition duration-100 ease-in"
                                leave-from-class="opacity-100"
                                leave-to-class="opacity-0"
                            >
                                <ListboxOptions
                                    class="absolute z-10 mt-1 max-h-28 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                >
                                    <ListboxOption
                                        v-for="theme in themes"
                                        :key="theme.key"
                                        v-slot="{ active, selected }"
                                        as="template"
                                        :value="theme.key"
                                    >
                                        <li
                                            :class="[
                                                active ? 'bg-primary-500' : '',
                                                'relative cursor-default select-none py-2 pl-8 pr-4',
                                            ]"
                                        >
                                            <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">
                                                {{ theme.name }}
                                            </span>

                                            <span
                                                v-if="selected"
                                                :class="[
                                                    active ? 'text-neutral' : 'text-primary-500',
                                                    'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                ]"
                                            >
                                                <CheckIcon class="size-5" />
                                            </span>
                                        </li>
                                    </ListboxOption>
                                </ListboxOptions>
                            </transition>
                        </div>
                    </Listbox>
                </template>
            </GenericContainerPanelEntry>
            <GenericContainerPanelEntry>
                <template #title>
                    {{ $t('components.auth.settings_panel.set_startpage.title') }}
                </template>
                <template #default>
                    <select
                        v-if="activeChar"
                        v-model="selectedHomepage"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    >
                        <option
                            v-for="page in homepages"
                            :key="page.path"
                            :value="page"
                            :disabled="!(page.permission === undefined || can(page.permission))"
                        >
                            {{ $t(page.name ?? 'common.page') }}
                        </option>
                    </select>
                    <p v-else>
                        {{ $t('components.auth.settings_panel.set_startpage.no_char_selected') }}
                    </p>
                </template>
            </GenericContainerPanelEntry>
            <GenericContainerPanelEntry>
                <template #title>
                    {{ $t('components.auth.settings_panel.editor_theme.title') }}
                </template>
                <template #default>
                    <div class="flex items-center">
                        <UToggle v-model="darkModeActive">
                            <span class="sr-only">{{ $t('components.auth.settings_panel.editor_theme.title') }}</span>
                        </UToggle>
                        <span class="ml-3 text-sm font-medium text-gray-300">{{
                            $t('components.auth.settings_panel.editor_theme.dark_mode')
                        }}</span>
                    </div>
                </template>
            </GenericContainerPanelEntry>
            <GenericContainerPanelEntry>
                <template #title>
                    {{ $t('components.auth.settings_panel.streamer_mode.title') }}
                </template>
                <template #default>
                    <div class="flex items-center">
                        <UToggle v-model="streamerMode">
                            <span class="sr-only">{{ $t('components.auth.settings_panel.streamer_mode.title') }}</span>
                        </UToggle>
                        <span class="ml-3 text-sm font-medium text-gray-300">{{
                            $t('components.auth.settings_panel.streamer_mode.description')
                        }}</span>
                    </div>
                </template>
            </GenericContainerPanelEntry>
        </template>
    </GenericContainerPanel>

    <GenericContainerPanel class="mt-3">
        <template #title>
            {{ $t('components.auth.settings_panel.volumes.title') }}
        </template>
        <template #description>
            {{ $t('components.auth.settings_panel.volumes.subtitle') }}
        </template>
        <template #default>
            <GenericContainerPanelEntry>
                <template #title>
                    {{ $t('common.notification', 2) }}
                </template>
                <template #default>
                    <label for="minmax-range" class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">
                        {{ $t('common.volume') }}:
                        {{ audio.notificationsVolume <= 0 ? 0 : (audio.notificationsVolume * 100).toFixed(0) }}%
                    </label>
                    <URange
                        v-model="audio.notificationsVolume"
                        :step="0.01"
                        :min="0"
                        :max="1"
                        class="h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 dark:bg-gray-700"
                    />
                </template>
            </GenericContainerPanelEntry>
        </template>
    </GenericContainerPanel>
</template>

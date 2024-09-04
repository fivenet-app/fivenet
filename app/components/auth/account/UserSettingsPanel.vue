<script lang="ts" setup>
import type { RoutePathSchema } from '@typed-router';
import { useAuthStore } from '~/store/auth';
import { useSettingsStore } from '~/store/settings';
import type { Perms } from '~~/gen/ts/perms';
import { backgroundColors, primaryColors } from './settings';

const { t } = useI18n();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const settings = useSettingsStore();
const { startpage, design, streamerMode, audio } = storeToRefs(settings);

const homepages: { name: string; path: RoutePathSchema; permission?: Perms }[] = [
    { name: t('common.overview'), path: '/overview' },
    { name: t('pages.citizens.title'), path: '/citizens', permission: 'CitizenStoreService.ListCitizens' },
    { name: t('pages.vehicles.title'), path: '/vehicles', permission: 'DMVService.ListVehicles' },
    { name: t('pages.documents.title'), path: '/documents', permission: 'DocStoreService.ListDocuments' },
    { name: t('pages.jobs.overview.title'), path: '/jobs/overview' },
    { name: t('common.calendar'), path: '/calendar' },
    { name: t('common.qualification', 2), path: '/qualifications', permission: 'QualificationsService.ListQualifications' },
    { name: t('common.livemap'), path: '/livemap', permission: 'LivemapperService.Stream' },
    { name: t('common.dispatch_center'), path: '/centrum', permission: 'CentrumService.TakeControl' },
];

const selectedHomepage = ref<(typeof homepages)[0]>();
watch(selectedHomepage, () => (startpage.value = selectedHomepage.value?.path ?? '/overview'));

onBeforeMount(async () => (selectedHomepage.value = homepages.find((h) => h.path === startpage.value)));

const darkModeActive = ref(design.value.documents.editorTheme === 'dark');

watch(darkModeActive, async () => {
    if (darkModeActive.value) {
        design.value.documents.editorTheme = 'dark';
    } else {
        design.value.documents.editorTheme = 'default';
    }
});

const designDocumentsListStyle = ref(design.value.documents.listStyle === 'double');

watch(designDocumentsListStyle, async () => {
    if (designDocumentsListStyle.value) {
        design.value.documents.listStyle = 'double';
    } else {
        design.value.documents.listStyle = 'single';
    }
});

const availableColorOptions = [...primaryColors, ...backgroundColors].map((color) => ({
    label: color,
    chip: color,
}));

const appConfig = useAppConfig();

watch(design.value, () => {
    appConfig.ui.primary = design.value.ui.primary;
    appConfig.ui.gray = design.value.ui.gray;
});
</script>

<template>
    <UDashboardPanelContent class="pb-24">
        <UDashboardSection :title="$t('common.theme')" :description="$t('components.auth.UserSettingsPanel.customization')">
            <template #links>
                <UColorModeSelect color="gray" />
            </template>

            <UFormGroup name="primaryColor" :label="$t('common.color')" class="grid grid-cols-2 items-center gap-2">
                <USelectMenu
                    v-model="design.ui.primary"
                    name="primaryColor"
                    :options="availableColorOptions"
                    option-attribute="label"
                    value-attribute="chip"
                    :searchable-placeholder="$t('common.search_field')"
                >
                    <template #label>
                        <span
                            class="size-2 rounded-full"
                            :class="`bg-${design.ui.primary}-500 dark:bg-${design.ui.primary}-400`"
                        />
                        <span class="truncate">{{ design.ui.primary }}</span>
                    </template>

                    <template #option="{ option }">
                        <span class="size-2 rounded-full" :class="`bg-${option.chip}-500 dark:bg-${option.chip}-400`" />
                        <span class="truncate">{{ option.label }}</span>
                    </template>
                </USelectMenu>
            </UFormGroup>

            <UFormGroup
                name="grayColor"
                :label="$t('components.auth.UserSettingsPanel.background_color')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <USelectMenu
                    v-model="design.ui.gray"
                    name="grayColor"
                    :options="availableColorOptions"
                    option-attribute="label"
                    value-attribute="chip"
                    :searchable-placeholder="$t('common.search_field')"
                >
                    <template #label>
                        <span class="size-2 rounded-full" :class="`bg-${design.ui.gray}-500 dark:bg-${design.ui.gray}-400`" />
                        <span class="truncate">{{ design.ui.gray }}</span>
                    </template>

                    <template #option="{ option }">
                        <span class="size-2 rounded-full" :class="`bg-${option.chip}-500 dark:bg-${option.chip}-400`" />
                        <span class="truncate">{{ option.label }}</span>
                    </template>
                </USelectMenu>
            </UFormGroup>

            <UFormGroup
                name="darkModeActive"
                :label="$t('components.auth.UserSettingsPanel.editor_theme.title')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <UToggle v-model="darkModeActive">
                    <span class="sr-only">{{ $t('components.auth.UserSettingsPanel.editor_theme.title') }}</span>
                </UToggle>
            </UFormGroup>

            <UFormGroup
                name="designDocumentsListStyle"
                :label="$t('components.auth.UserSettingsPanel.documents_lists_style.title')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <div class="inline-flex items-center gap-2">
                    <span>{{ $t('components.auth.UserSettingsPanel.documents_lists_style.single') }}</span>
                    <UToggle v-model="designDocumentsListStyle">
                        <span class="sr-only">{{ $t('components.auth.UserSettingsPanel.documents_lists_style.title') }}</span>
                    </UToggle>
                    <span>{{ $t('components.auth.UserSettingsPanel.documents_lists_style.double') }}</span>
                </div>
            </UFormGroup>

            <UFormGroup
                name="streamerMode"
                :label="$t('components.auth.UserSettingsPanel.streamer_mode.title')"
                :description="$t('components.auth.UserSettingsPanel.streamer_mode.description')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <UToggle v-model="streamerMode" name="streamerMode">
                    <span class="sr-only">{{ $t('components.auth.UserSettingsPanel.streamer_mode.title') }}</span>
                </UToggle>
            </UFormGroup>

            <UFormGroup
                name="selectedHomepage"
                :label="$t('components.auth.UserSettingsPanel.set_startpage.title')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <USelectMenu
                    v-if="activeChar"
                    v-model="selectedHomepage"
                    :options="homepages.filter((h) => h.permission === undefined || can(h.permission).value)"
                    option-attribute="name"
                    :searchable-placeholder="$t('common.search_field')"
                />
                <p v-else class="text-sm">
                    {{ $t('components.auth.UserSettingsPanel.set_startpage.no_char_selected') }}
                </p>
            </UFormGroup>
        </UDashboardSection>

        <UDivider class="mb-4" />

        <UDashboardSection
            :title="$t('components.auth.UserSettingsPanel.volumes.title')"
            :description="$t('components.auth.UserSettingsPanel.volumes.subtitle')"
        >
            <UFormGroup
                name="notificationsVolume"
                :label="$t('components.auth.UserSettingsPanel.volumes.notifications_volume')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <URange v-model="audio.notificationsVolume" :step="0.01" :min="0" :max="1" />
                {{ audio.notificationsVolume <= 0 ? 0 : (audio.notificationsVolume * 100).toFixed(0) }}%
            </UFormGroup>
        </UDashboardSection>
    </UDashboardPanelContent>
</template>

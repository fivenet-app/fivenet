<script lang="ts" setup>
import type { RoutePathSchema } from '@typed-router';
import { useAuthStore } from '~/store/auth';
import { useSettingsStore } from '~/store/settings';
import type { Perms } from '~~/gen/ts/perms';

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
    { name: t('pages.jobs.overview.title'), path: '/jobs/overview', permission: 'JobsService.ListColleagues' },
    { name: t('common.livemap'), path: '/livemap', permission: 'LivemapperService.Stream' },
    { name: t('common.dispatch_center'), path: '/centrum', permission: 'CentrumService.TakeControl' },
];

const selectedHomepage = ref<(typeof homepages)[0]>();
watch(selectedHomepage, () => (startpage.value = selectedHomepage.value?.path ?? '/overview'));

onBeforeMount(async () => {
    selectedHomepage.value = homepages.find((h) => h.path === startpage.value);
});

const darkModeActive = ref(design.value.docEditorTheme === 'dark');

watch(darkModeActive, async () => {
    if (darkModeActive.value) {
        design.value.docEditorTheme = 'dark';
    } else {
        design.value.docEditorTheme = 'default';
    }
});

const colors = ['green', 'teal', 'cyan', 'sky', 'blue', 'indigo', 'violet'];

const appConfig = useAppConfig();

watch(design, () => (appConfig.ui.primary = design.value.ui.primary));
</script>

<template>
    <UDashboardPanelContent class="pb-24">
        <UDashboardSection :title="$t('common.theme')" description="Customize the look and feel of your dashboard.">
            <template #links>
                <UColorModeSelect color="gray" />
            </template>

            <UFormGroup name="primaryColor" :label="$t('common.color')" class="grid grid-cols-2 items-center gap-2">
                <USelectMenu v-model="design.ui.primary" :options="colors">
                    <template #label>
                        <span
                            class="h-2 w-2 rounded-full"
                            :class="`bg-${design.ui.primary}-500 dark:bg-${design.ui.primary}-400`"
                        />
                        <span class="truncate">{{ design.ui.primary }}</span>
                    </template>
                    <template #option="{ option }">
                        <span class="h-2 w-2 rounded-full" :class="`bg-${option}-500 dark:bg-${option}-400`" />
                        <span class="truncate">{{ option }}</span>
                    </template>
                </USelectMenu>
            </UFormGroup>

            <UFormGroup
                name="darkModeActive"
                :label="$t('components.auth.settings_panel.editor_theme.title')"
                class="grid grid-cols-2 items-center gap-2"
                :ui="{ container: 'justify-self-end' }"
            >
                <UToggle v-model="darkModeActive">
                    <span class="sr-only">{{ $t('components.auth.settings_panel.editor_theme.title') }}</span>
                </UToggle>
            </UFormGroup>
        </UDashboardSection>

        <UDivider class="mb-4" />

        <UDashboardSection
            :title="$t('components.auth.settings_panel.streamer_mode.title')"
            :description="$t('components.auth.settings_panel.streamer_mode.description')"
        >
            <template #links>
                <UToggle v-model="streamerMode">
                    <span class="sr-only">{{ $t('components.auth.settings_panel.streamer_mode.title') }}</span>
                </UToggle>
            </template>
        </UDashboardSection>

        <UDivider class="mb-4" />

        <UDashboardSection :title="$t('components.auth.settings_panel.title')">
            <UFormGroup
                name="selectedHomepage"
                :label="$t('components.auth.settings_panel.set_startpage.title')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <USelectMenu
                    v-if="activeChar"
                    v-model="selectedHomepage"
                    :options="homepages.filter((h) => h.permission === undefined || can(h.permission))"
                    option-attribute="name"
                />
                <p v-else class="text-sm">
                    {{ $t('components.auth.settings_panel.set_startpage.no_char_selected') }}
                </p>
            </UFormGroup>
        </UDashboardSection>

        <UDivider class="mb-4" />

        <UDashboardSection
            :title="$t('components.auth.settings_panel.volumes.title')"
            :description="$t('components.auth.settings_panel.volumes.subtitle')"
        >
            <UFormGroup
                name="selectedHomepage"
                :label="$t('components.auth.settings_panel.set_startpage.title')"
                class="grid grid-cols-2 items-center gap-2"
            >
                <URange v-model="audio.notificationsVolume" :step="0.01" :min="0" :max="1" />
                {{ audio.notificationsVolume <= 0 ? 0 : (audio.notificationsVolume * 100).toFixed(0) }}%
            </UFormGroup>
        </UDashboardSection>

        <UDivider class="mb-4" />
    </UDashboardPanelContent>
</template>

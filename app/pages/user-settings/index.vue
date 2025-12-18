<script lang="ts" setup>
import type { RoutePathSchema } from '@typed-router';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const { can, activeChar } = useAuth();

const settingsStore = useSettingsStore();
const { startpage, design, streamerMode, eventsDisabled } = storeToRefs(settingsStore);

const startpages: { label: string; path: RoutePathSchema; permission?: Perms }[] = [
    { label: t('common.overview'), path: '/overview' },
    { label: t('common.mail'), path: '/mail/:thread?', permission: 'mailer.MailerService/ListEmails' },
    { label: t('pages.citizens.title'), path: '/citizens', permission: 'citizens.CitizensService/ListCitizens' },
    { label: t('pages.vehicles.title'), path: '/vehicles', permission: 'vehicles.VehiclesService/ListVehicles' },
    { label: t('pages.documents.title'), path: '/documents/', permission: 'documents.DocumentsService/ListDocuments' },
    { label: t('pages.jobs.overview.title'), path: '/jobs/overview', permission: 'jobs.JobsService/ListColleagues' },
    { label: t('common.calendar'), path: '/calendar' },
    {
        label: t('common.qualification', 2),
        path: '/qualifications',
        permission: 'qualifications.QualificationsService/ListQualifications',
    },
    { label: t('common.livemap'), path: '/livemap', permission: 'livemap.LivemapService/Stream' },
    { label: t('common.dispatch_center'), path: '/centrum', permission: 'centrum.CentrumService/TakeControl' },
    { label: t('common.wiki'), path: '/wiki', permission: 'wiki.WikiService/ListPages' },
];

const selectedHomepage = ref<(typeof startpages)[0]>();
watch(selectedHomepage, () => (startpage.value = selectedHomepage.value?.path ?? '/overview'));

const designDocumentListStyle = ref(design.value.documents.listStyle === 'double');

watch(designDocumentListStyle, async () => {
    if (designDocumentListStyle.value) {
        design.value.documents.listStyle = 'double';
    } else {
        design.value.documents.listStyle = 'single';
    }
});
</script>

<template>
    <UPageCard :title="$t('common.theme')" :description="$t('components.auth.user_settings_panel.customization')">
        <template #links>
            <UColorModeSelect color="neutral" />
        </template>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="primaryColor" :label="$t('common.color')">
            <ColorPickerTW v-model="design.ui.primary" name="primaryColor" class="w-full" />
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="grayColor"
            :label="$t('components.auth.user_settings_panel.background_color')"
        >
            <ColorPickerTW v-model="design.ui.gray" name="grayColor" class="w-full" />
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="streamerMode"
            :label="$t('components.auth.user_settings_panel.streamer_mode.title')"
            :description="$t('components.auth.user_settings_panel.streamer_mode.description')"
        >
            <USwitch v-model="streamerMode" name="streamerMode" />
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="designDocumentListStyle"
            :label="$t('components.auth.user_settings_panel.documents_lists_style.title')"
        >
            <div class="inline-flex items-center gap-2 text-sm">
                <span>{{ $t('components.auth.user_settings_panel.documents_lists_style.single') }}</span>
                <USwitch v-model="designDocumentListStyle" />
                <span>{{ $t('components.auth.user_settings_panel.documents_lists_style.double') }}</span>
            </div>
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="selectedHomepage"
            :label="$t('components.auth.user_settings_panel.set_startpage.title')"
        >
            <ClientOnly v-if="activeChar">
                <USelectMenu
                    v-model="selectedHomepage"
                    :items="startpages.filter((h) => h.permission === undefined || can(h.permission).value)"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    class="w-full"
                />
            </ClientOnly>
            <p v-else class="text-sm">
                {{ $t('components.auth.user_settings_panel.set_startpage.no_char_selected') }}
            </p>
        </UFormField>

        <UFormField
            :label="$t('components.auth.user_settings_panel.events_disabled.label')"
            :description="$t('components.auth.user_settings_panel.events_disabled.description')"
            name="eventsDisabled"
            class="grid grid-cols-2 items-center gap-2"
        >
            <USwitch v-model="eventsDisabled" />
        </UFormField>
    </UPageCard>
</template>

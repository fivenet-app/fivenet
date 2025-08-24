<script lang="ts" setup>
import type { RoutePathSchema } from '@typed-router';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const { can, activeChar } = useAuth();

const settings = useSettingsStore();
const { startpage, design, streamerMode } = storeToRefs(settings);

const startpages: { name: string; path: RoutePathSchema; permission?: Perms }[] = [
    { name: t('common.overview'), path: '/overview' },
    { name: t('common.mail'), path: '/mail', permission: 'mailer.MailerService/ListEmails' },
    { name: t('pages.citizens.title'), path: '/citizens', permission: 'citizens.CitizensService/ListCitizens' },
    { name: t('pages.vehicles.title'), path: '/vehicles', permission: 'vehicles.VehiclesService/ListVehicles' },
    { name: t('pages.documents.title'), path: '/documents', permission: 'documents.DocumentsService/ListDocuments' },
    { name: t('pages.jobs.overview.title'), path: '/jobs/overview', permission: 'jobs.JobsService/ListColleagues' },
    { name: t('common.calendar'), path: '/calendar' },
    {
        name: t('common.qualification', 2),
        path: '/qualifications',
        permission: 'qualifications.QualificationsService/ListQualifications',
    },
    { name: t('common.livemap'), path: '/livemap', permission: 'livemap.LivemapService/Stream' },
    { name: t('common.dispatch_center'), path: '/centrum', permission: 'centrum.CentrumService/TakeControl' },
    { name: t('common.wiki'), path: '/wiki', permission: 'wiki.WikiService/ListPages' },
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
    <UPageCard :title="$t('common.theme')" :description="$t('components.auth.UserSettingsPanel.customization')">
        <template #links>
            <UColorModeSelect color="neutral" />
        </template>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="primaryColor" :label="$t('common.color')">
            <ColorPickerTW v-model="design.ui.primary" name="primaryColor" />
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="grayColor"
            :label="$t('components.auth.UserSettingsPanel.background_color')"
        >
            <ColorPickerTW v-model="design.ui.gray" name="grayColor" />
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="streamerMode"
            :label="$t('components.auth.UserSettingsPanel.streamer_mode.title')"
            :description="$t('components.auth.UserSettingsPanel.streamer_mode.description')"
        >
            <USwitch v-model="streamerMode" name="streamerMode" />
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="designDocumentListStyle"
            :label="$t('components.auth.UserSettingsPanel.documents_lists_style.title')"
        >
            <div class="inline-flex items-center gap-2 text-sm">
                <span>{{ $t('components.auth.UserSettingsPanel.documents_lists_style.single') }}</span>
                <USwitch v-model="designDocumentListStyle" />
                <span>{{ $t('components.auth.UserSettingsPanel.documents_lists_style.double') }}</span>
            </div>
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="selectedHomepage"
            :label="$t('components.auth.UserSettingsPanel.set_startpage.title')"
        >
            <ClientOnly v-if="activeChar">
                <USelectMenu
                    v-model="selectedHomepage"
                    :items="startpages.filter((h) => h.permission === undefined || can(h.permission).value)"
                    option-attribute="name"
                    :searchable-placeholder="$t('common.search_field')"
                />
            </ClientOnly>
            <p v-else class="text-sm">
                {{ $t('components.auth.UserSettingsPanel.set_startpage.no_char_selected') }}
            </p>
        </UFormField>
    </UPageCard>
</template>

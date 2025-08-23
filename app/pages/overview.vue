<script lang="ts" setup>
import FiveNetHints from '~/components/FiveNetHints.vue';
import CardsList from '~/components/partials/CardsList.vue';
import type { CardElements } from '~/utils/types';

useHead({
    title: 'common.overview',
});

definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
});

const { t } = useI18n();

const items = [
    {
        title: t('common.mail'),
        description: t('pages.overview.features.mailer'),
        to: { name: 'mail' },
        permission: 'mailer.MailerService/ListEmails',
        icon: 'i-mdi-inbox-full-outline',
    },
    {
        title: `${t('common.citizen', 2)} ${t('common.search')}`,
        description: t('pages.overview.features.citizens'),
        to: { name: 'citizens' },
        permission: 'citizens.CitizensService/ListCitizens',
        icon: 'i-mdi-account-multiple-outline',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        to: { name: 'vehicles' },
        permission: 'vehicles.VehiclesService/ListVehicles',
        icon: 'i-mdi-car-outline',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        to: { name: 'documents' },
        permission: 'documents.DocumentsService/ListDocuments',
        icon: 'i-mdi-file-document-box-multiple-outline',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        to: { name: 'jobs-overview' },
        permission: 'jobs.JobsService/ListColleagues',
        icon: 'i-mdi-briefcase-outline',
    },
    {
        title: t('common.calendar'),
        description: t('pages.overview.features.calendar'),
        to: { name: 'calendar' },
        icon: 'i-mdi-calendar-outline',
    },
    {
        title: t('common.qualification', 2),
        description: t('pages.overview.features.qualifications'),
        to: { name: 'qualifications' },
        permission: 'qualifications.QualificationsService/ListQualifications',
        icon: 'i-mdi-school-outline',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        to: { name: 'livemap' },
        permission: 'livemap.LivemapService/Stream',
        icon: 'i-mdi-map-outline',
    },
    {
        title: t('common.dispatch_center'),
        description: t('pages.overview.features.centrum'),
        to: { name: 'centrum' },
        permission: 'centrum.CentrumService/TakeControl',
        icon: 'i-mdi-car-emergency',
    },
    {
        title: t('common.wiki'),
        description: t('pages.overview.features.wiki'),
        to: { name: 'wiki' },
        permission: 'wiki.WikiService/ListPages',
        icon: 'i-mdi-brain',
    },
] as CardElements;
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('common.overview')" />
        </template>

        <template #body>
            <CardsList :items="items" />

            <FiveNetHints class="mt-4" />
        </template>
    </UDashboardPanel>
</template>

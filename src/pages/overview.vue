<script lang="ts" setup>
import FiveNetHints from '~/components/FiveNetHints.vue';
import CardsList from '~/components/partials/CardsList.vue';
import { type CardElements } from '~/utils/types';

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
        title: t('common.messenger'),
        description: t('pages.overview.features.messenger'),
        to: { name: 'messenger' },
        permission: 'MessengerService.ListThreads',
        icon: 'i-mdi-conversation-outline',
    },
    {
        title: `${t('common.citizen', 2)} ${t('common.search')}`,
        description: t('pages.overview.features.citizens'),
        to: { name: 'citizens' },
        permission: 'CitizenStoreService.ListCitizens',
        icon: 'i-mdi-account-multiple-outline',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        to: { name: 'vehicles' },
        permission: 'DMVService.ListVehicles',
        icon: 'i-mdi-car-outline',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        to: { name: 'documents' },
        permission: 'DocStoreService.ListDocuments',
        icon: 'i-mdi-file-document-box-multiple-outline',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        to: { name: 'jobs-overview' },
        permission: 'JobsService.ListColleagues',
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
        permission: 'QualificationsService.ListQualifications',
        icon: 'i-mdi-school-outline',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        to: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: 'i-mdi-map-outline',
    },
    {
        title: t('common.dispatch_center'),
        description: t('pages.overview.features.centrum'),
        to: { name: 'centrum' },
        permission: 'CentrumService.TakeControl',
        icon: 'i-mdi-car-emergency',
    },
] as CardElements;
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.overview')" />

            <UDashboardPanelContent>
                <CardsList :items="items" />

                <FiveNetHints class="mt-4" />
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>

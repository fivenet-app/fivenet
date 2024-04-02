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

const features = [
    {
        title: `${t('common.citizen', 2)} ${t('common.search')}`,
        description: t('pages.overview.features.citizens'),
        to: { name: 'citizens' },
        permission: 'CitizenStoreService.ListCitizens',
        icon: 'i-mdi-account-multiple-outline',
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-100',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        to: { name: 'vehicles' },
        permission: 'DMVService.ListVehicles',
        icon: 'i-mdi-car-outline',
        iconForeground: 'text-zinc-900',
        iconBackground: 'bg-zinc-100',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        to: { name: 'documents' },
        permission: 'DocStoreService.ListDocuments',
        icon: 'i-mdi-file-document-box-multiple-outline',
        iconForeground: 'text-sky-900',
        iconBackground: 'bg-sky-100',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        to: { name: 'jobs-overview' },
        permission: 'JobsService.ListColleagues',
        icon: 'i-mdi-briefcase-outline',
        iconForeground: 'text-warn-900',
        iconBackground: 'bg-warn-100',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        to: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: 'i-mdi-map-outline',
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-100',
    },
    {
        title: t('common.dispatch_center'),
        description: t('pages.overview.features.centrum'),
        to: { name: 'centrum' },
        permission: 'CentrumService.TakeControl',
        icon: 'i-mdi-car-emergency',
        iconForeground: 'text-error-900',
        iconBackground: 'bg-error-100',
    },
] as CardElements;
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.overview')"> </UDashboardNavbar>

            <CardsList class="mt-2" :items="features" />

            <FiveNetHints class="m-2" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

<script lang="ts" setup>
import { mdiAccountMultiple, mdiBriefcase, mdiCar, mdiCarEmergency, mdiFileDocumentMultiple, mdiMap } from '@mdi/js';
import Cards from '~/components/partials/Cards.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import { CardElements } from '~/utils/types';

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
        href: { name: 'citizens' },
        permission: 'CitizenStoreService.ListCitizens',
        icon: mdiAccountMultiple,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-100',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        href: { name: 'vehicles' },
        permission: 'DMVService.ListVehicles',
        icon: mdiCar,
        iconForeground: 'text-zinc-900',
        iconBackground: 'bg-zinc-100',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        href: { name: 'documents' },
        permission: 'DocStoreService.ListDocuments',
        icon: mdiFileDocumentMultiple,
        iconForeground: 'text-sky-900',
        iconBackground: 'bg-sky-100',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        href: { name: 'jobs' },
        permission: 'Jobs.View',
        icon: mdiBriefcase,
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-100',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        href: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: mdiMap,
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-100',
    },
    {
        title: t('common.dispatch_center'),
        description: t('pages.overview.features.livemap'),
        href: { name: 'centrum' },
        permission: 'CentrumService.Stream',
        icon: mdiCarEmergency,
        iconForeground: 'text-red-900',
        iconBackground: 'bg-red-100',
    },
] as CardElements;
</script>

<template>
    <ContentCenterWrapper>
        <Cards :items="features" />
    </ContentCenterWrapper>
</template>

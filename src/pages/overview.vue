<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import { BriefcaseIcon, DocumentTextIcon, UsersIcon, MapIcon, TruckIcon } from '@heroicons/vue/24/outline';
import Cards from '~/components/partials/Cards.vue';
import { CardElements } from '~/utils/types';

useHead({
    title: 'Overview',
});
definePageMeta({
    title: 'Overview',
    requiresAuth: true,
});

const { t } = useI18n();

const features = [
    {
        title: t('common.citizen', 2),
        description: t('pages.overview.features.citizens'),
        href: { name: 'citizens' },
        permission: 'CitizenStoreService.FindUsers',
        icon: UsersIcon,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-50',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        href: { name: 'vehicles' },
        permission: 'DMVService.FindVehicles',
        icon: TruckIcon,
        iconForeground: 'text-zinc-900',
        iconBackground: 'bg-zinc-50',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        href: { name: 'documents' },
        permission: 'DocStoreService.FindDocuments',
        icon: DocumentTextIcon,
        iconForeground: 'text-sky-900',
        iconBackground: 'bg-sky-50',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        href: { name: 'jobs' },
        permission: 'Jobs.View',
        icon: BriefcaseIcon,
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-50',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        href: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: MapIcon,
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-50',
    },
] as CardElements;
</script>

<template>
    <ContentCenterWrapper>
        <Cards :items="features" />
    </ContentCenterWrapper>
</template>

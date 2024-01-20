<script lang="ts" setup>
import { AccountMultipleIcon, BriefcaseIcon, CarEmergencyIcon, CarIcon, FileDocumentMultipleIcon, MapIcon } from 'mdi-vue3';
import FiveNetHints from '~/components/FiveNetHints.vue';
import CardsList from '~/components/partials/CardsList.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
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
        href: { name: 'citizens' },
        permission: 'CitizenStoreService.ListCitizens',
        icon: markRaw(AccountMultipleIcon),
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-100',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        href: { name: 'vehicles' },
        permission: 'DMVService.ListVehicles',
        icon: markRaw(CarIcon),
        iconForeground: 'text-zinc-900',
        iconBackground: 'bg-zinc-100',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        href: { name: 'documents' },
        permission: 'DocStoreService.ListDocuments',
        icon: markRaw(FileDocumentMultipleIcon),
        iconForeground: 'text-sky-900',
        iconBackground: 'bg-sky-100',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        href: { name: 'jobs-overview' },
        permission: 'JobsService.ListColleagues',
        icon: markRaw(BriefcaseIcon),
        iconForeground: 'text-warn-900',
        iconBackground: 'bg-warn-100',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        href: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: markRaw(MapIcon),
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-100',
    },
    {
        title: t('common.dispatch_center'),
        description: t('pages.overview.features.centrum'),
        href: { name: 'centrum' },
        permission: 'CentrumService.TakeControl',
        icon: markRaw(CarEmergencyIcon),
        iconForeground: 'text-error-900',
        iconBackground: 'bg-error-100',
    },
] as CardElements;
</script>

<template>
    <ContentCenterWrapper>
        <CardsList :items="features" />
        <FiveNetHints />
    </ContentCenterWrapper>
</template>

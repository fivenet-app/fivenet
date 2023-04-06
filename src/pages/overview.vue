<script lang="ts" setup>
import ContentCenterWrapper from '../components/partials/ContentCenterWrapper.vue';
import { BriefcaseIcon, DocumentTextIcon, UsersIcon, MapIcon, TruckIcon } from '@heroicons/vue/24/outline';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import Cards from '../components/partials/Cards.vue';
import { CardElements } from '../utils/types';

useHead({
    title: 'Overview',
});
definePageMeta({
    title: 'Overview',
    requiresAuth: true,
});

const features = [
    {
        title: 'Citizen Search',
        description:
            'Search and find information about Citizens, including basic info, their licenses and related documents.',
        href: { name: 'citizens' },
        permission: 'CitizenStoreService.FindUsers',
        icon: UsersIcon,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-50',
    },
    {
        title: 'Vehicles',
        description:
            'Search and find information about Vehicles.',
        href: { name: 'vehicles' },
        permission: 'DMVService.FindVehicles',
        icon: TruckIcon,
        iconForeground: 'text-zinc-900',
        iconBackground: 'bg-zinc-50',
    },
    {
        title: 'Documents',
        description:
            'Search and find information about Citizens, including basic info, their licenses and related documents.',
        href: { name: 'documents' },
        permission: 'DocStoreService.FindDocuments',
        icon: DocumentTextIcon,
        iconForeground: 'text-sky-900',
        iconBackground: 'bg-sky-50',
    },
    {
        title: 'Job',
        description: 'Infos about your job and employee management. Coming soon',
        href: { name: 'jobs' },
        permission: 'Jobs.View',
        icon: BriefcaseIcon,
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-50',
    },
    {
        title: 'Livemap',
        description: 'Live position of dispatches and your colleagues.',
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

<script setup lang="ts">
import { BriefcaseIcon, DocumentTextIcon, UsersIcon, MapIcon, TruckIcon } from '@heroicons/vue/24/outline';
import { ArrowUpRightIcon } from '@heroicons/vue/24/solid';

const features = [
    {
        title: 'Citizen Search',
        description:
            'Search and find information about Citizens, including basic info, their licenses and related documents.',
        href: '/citizens',
        permission: 'CitizenStoreService.FindUsers',
        icon: UsersIcon,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-50',
    },
    {
        title: 'Vehicles',
        description:
            'Search and find information about Vehicles.',
        href: '/vehicles',
        permission: 'DMVService.FindVehicles',
        icon: TruckIcon,
        iconForeground: 'text-zinc-900',
        iconBackground: 'bg-zinc-50',
    },
    {
        title: 'Documents',
        description:
            'Search and find information about Citizens, including basic info, their licenses and related documents.',
        href: '/documents',
        permission: 'DocStoreService.FindDocuments',
        icon: DocumentTextIcon,
        iconForeground: 'text-sky-900',
        iconBackground: 'bg-sky-50',
    },
    {
        title: 'Job',
        description: 'Infos about your job and employee management. Coming soon',
        href: '/job',
        permission: 'Jobs.View',
        icon: BriefcaseIcon,
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-50',
    },
    {
        title: 'Livemap',
        description: 'Live position of dispatches and your colleagues.',
        href: '/livemap',
        permission: 'LivemapperService.Stream',
        icon: MapIcon,
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-50',
    },
]
</script>

<template>
    <div
        class="overflow-hidden divide-y-4 rounded-lg bg-base-900 shadow-float sm:grid sm:grid-cols-2 sm:gap-1 sm:max-w-6xl sm:mx-auto divide-base-900 sm:divide-y-0">
        <div v-for="(feature, featureIdx) in features" v-can="feature.permission" :key="feature.title" :class="[
            featureIdx === 0 ? 'rounded-tl-lg rounded-tr-lg sm:rounded-tr-none' : '',
            featureIdx === 1 ? 'sm:rounded-tr-lg' : '',
            featureIdx === features.length - 2 && features.length % 2 === 0 ? 'sm:rounded-bl-lg' : '',
            featureIdx === features.length - 1 && features.length % 2 === 0 ? 'rounded-br-lg' : '',
            featureIdx === features.length - 1 ? 'rounded-bl-lg sm:rounded-bl-none' : '',
            'group relative bg-base-700 p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-neutral',
        ]">
            <div>
                <span
                    :class="[feature.iconBackground, feature.iconForeground, 'inline-flex rounded-lg p-3']">
                    <component :is="feature.icon" class="h-auto w-7" aria-hidden="true" />
                </span>
            </div>
            <div class="mt-4">
                <h3 class="text-base font-semibold leading-6 text-neutral">
                    <router-link :to="feature.href" class="focus:outline-none">
                        <!-- Extend touch target to entire panel -->
                        <span class="absolute inset-0" aria-hidden="true" />
                        {{ feature.title }}
                    </router-link>
                </h3>
                <p class="mt-2 text-sm text-base-200">{{ feature.description }}</p>
            </div>
            <span class="absolute pointer-events-none top-6 right-6 text-base-300 group-hover:text-base-200"
                aria-hidden="true">
                <ArrowUpRightIcon class="w-6 h-6" />
            </span>
        </div>
    </div>
</template>

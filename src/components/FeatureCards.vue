<script lang="ts">
import { defineComponent } from 'vue';
import { BriefcaseIcon, BellAlertIcon, DocumentTextIcon, UsersIcon, MapIcon } from '@heroicons/vue/24/outline';

export default defineComponent({
    data() {
        return {
            features: [
                {
                    title: 'Citizen Search',
                    description:
                        'Search and find information about Citizens, including basic info, their licenses and related documents.',
                    href: '/citizens',
                    permission: 'users-findusers',
                    icon: UsersIcon,
                    iconForeground: 'text-purple-800',
                    iconBackground: 'bg-purple-50',
                },
                {
                    title: 'Documents',
                    description:
                        'Search and find information about Citizens, including basic info, their licenses and related documents.',
                    href: '/documents',
                    permission: 'documents-view',
                    icon: DocumentTextIcon,
                    iconForeground: 'text-sky-800',
                    iconBackground: 'bg-sky-50',
                },
                {
                    title: 'Dispatches',
                    description: 'Dispatch Center and Management. Coming soon.',
                    href: '/dispatches',
                    permission: 'dispatches-view',
                    icon: BellAlertIcon,
                    iconForeground: 'text-rose-800',
                    iconBackground: 'bg-rose-50',
                },
                {
                    title: 'Job',
                    description: 'Infos about your job and employee management. Coming soon',
                    href: '/job',
                    permission: 'job-view',
                    icon: BriefcaseIcon,
                    iconForeground: 'text-yellow-800',
                    iconBackground: 'bg-yellow-50',
                },
                {
                    title: 'Livemap',
                    description: 'Live position of dispatches and your colleagues.',
                    href: '/livemap',
                    permission: 'livemap-stream',
                    icon: MapIcon,
                    iconForeground: 'text-teal-800',
                    iconBackground: 'bg-teal-50',
                },
            ],
        };
    },
});
</script>

<template>
    <div
        class="divide-y divide-white overflow-hidden rounded-lg bg-gray-800 shadow sm:grid sm:grid-cols-2 sm:gap-px sm:divide-y-0"
    >
        <div
            v-for="(feature, featureIdx) in features"
            v-can="feature.permission"
            :key="feature.title"
            :class="[
                featureIdx === 0 ? 'rounded-tl-lg rounded-tr-lg sm:rounded-tr-none' : '',
                featureIdx === 1 ? 'sm:rounded-tr-lg' : '',
                featureIdx === features.length - 2 ? 'sm:rounded-bl-lg' : '',
                featureIdx === features.length - 1 ? 'rounded-bl-lg rounded-br-lg sm:rounded-bl-none' : '',
                'group relative bg-gray-700 p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-white',
            ]"
        >
            <div>
                <span :class="[feature.iconBackground, feature.iconForeground, 'inline-flex rounded-lg p-3 ring-4 ring-white']">
                    <component :is="feature.icon" class="h-6 w-6" aria-hidden="true" />
                </span>
            </div>
            <div class="mt-8">
                <h3 class="text-base font-semibold leading-6 text-white">
                    <router-link :to="feature.href" class="focus:outline-none">
                        <!-- Extend touch target to entire panel -->
                        <span class="absolute inset-0" aria-hidden="true" />
                        {{ feature.title }}
                    </router-link>
                </h3>
                <p class="mt-2 text-sm text-gray-300">{{ feature.description }}</p>
            </div>
            <span class="pointer-events-none absolute top-6 right-6 text-gray-300 group-hover:text-gray-400" aria-hidden="true">
                <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                    <path
                        d="M20 4h1a1 1 0 00-1-1v1zm-1 12a1 1 0 102 0h-2zM8 3a1 1 0 000 2V3zM3.293 19.293a1 1 0 101.414 1.414l-1.414-1.414zM19 4v12h2V4h-2zm1-1H8v2h12V3zm-.707.293l-16 16 1.414 1.414 16-16-1.414-1.414z"
                    />
                </svg>
            </span>
        </div>
    </div>
</template>

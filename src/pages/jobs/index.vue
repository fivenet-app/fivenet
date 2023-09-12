<script lang="ts" setup>
import { CloseIcon, MenuIcon } from 'mdi-vue3';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';

const navigation: { name: string; to: RoutesNamedLocations; permission?: string }[] = [
    { name: 'common.overview', to: { name: 'jobs-index-overview' } },
    { name: 'pages.jobs.colleagues.title', to: { name: 'jobs-index-colleagues' } },
    { name: 'pages.jobs.requests.title', to: { name: 'jobs-index-requests' } },
    { name: 'pages.jobs.trainings.title', to: { name: 'jobs-index-trainings' } },
    //{ name: 'pages.jobs.timeclock.title', to: { name: 'jobs-index-timeclock' } },
    { name: 'pages.jobs.conduct.title', to: { name: 'jobs-index-conduct' }, permission: 'Jobs.ConductListEntries' },
];

useHead({
    title: 'pages.jobs.title',
});
definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: 'Jobs.ColleaguesList',
});

const route = useRoute();
const open = ref(false);
</script>

<template>
    <div class="min-h-full">
        <nav class="bg-primary-600">
            <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="flex h-16 items-center justify-between">
                    <div class="flex items-center">
                        <div class="-ml-2 flex md:hidden">
                            <!-- Mobile menu button -->
                            <button
                                type="button"
                                class="relative inline-flex items-center justify-center rounded-md bg-primary-600 p-2 text-primary-200 hover:bg-primary-500 hover:bg-opacity-75 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-primary-600"
                                @click="open = !open"
                            >
                                <span class="absolute -inset-0.5" />
                                <span class="sr-only">{{ $t('components.partials.sidebar.open_navigation') }}</span>
                                <MenuIcon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
                                <CloseIcon v-else class="block h-6 w-6" aria-hidden="true" />
                            </button>
                        </div>
                        <div class="hidden md:block">
                            <div class="flex items-baseline space-x-4">
                                <template v-for="item in navigation" :key="item.name">
                                    <NuxtLink
                                        v-if="item.permission === undefined || can(item.permission)"
                                        :to="item.to"
                                        class="text-white hover:bg-primary-500 hover:bg-opacity-75 rounded-md px-3 py-2 text-sm font-medium"
                                        active-class="bg-primary-700 text-white"
                                        aria-current-value="page"
                                    >
                                        {{ $t(item.name) }}
                                    </NuxtLink>
                                </template>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="md:hidden" :class="open ? 'block' : 'hidden'">
                <div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
                    <template v-for="item in navigation" :key="item.name">
                        <NuxtLink
                            v-if="item.permission === undefined || can(item.permission)"
                            :to="item.to"
                            class="text-white hover:bg-primary-500 hover:bg-opacity-75 block rounded-md px-3 py-2 text-base font-medium"
                            active-class="bg-primary-700 text-white"
                            aria-current-value="page"
                        >
                            {{ $t(item.name) }}
                        </NuxtLink>
                    </template>
                </div>
            </div>
        </nav>

        <header class="bg-base-700 shadow">
            <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
                <h1 class="text-3xl font-bold leading-tight tracking-tight text-white">
                    {{ $t(route.meta.title ?? 'pages.jobs.title') }}
                </h1>
            </div>
        </header>
        <main>
            <div class="mx-auto max-w-8xl py-6 sm:px-6 lg:px-8">
                <NuxtPage />
            </div>
        </main>
    </div>
</template>

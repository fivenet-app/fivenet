<script lang="ts" setup>
import { CloseIcon, MenuIcon } from 'mdi-vue3';
import { type RoutesNamedLocations } from '@typed-router';

const navigation: { name: string; to: RoutesNamedLocations; permission?: string }[] = [
    { name: 'common.overview', to: { name: 'jobs-overview' }, permission: 'JobsService.ListColleagues' },
    { name: 'pages.jobs.colleagues.title', to: { name: 'jobs-colleagues' }, permission: 'JobsService.ListColleagues' },
    {
        name: 'pages.jobs.requests.title',
        to: { name: 'jobs-requests' },
        permission: 'JobsRequestsService.ListRequestsRequestEntries',
    },
    {
        name: 'pages.jobs.qualifications.title',
        to: { name: 'jobs-qualifications' },
        permission: 'JobsService.QualificationsListEntries',
    },
    { name: 'pages.jobs.timeclock.title', to: { name: 'jobs-timeclock' }, permission: 'JobsTimeclockService.ListTimeclock' },
    { name: 'pages.jobs.conduct.title', to: { name: 'jobs-conduct' }, permission: 'JobsConductService.ListConductEntries' },
];

useHead({
    title: 'pages.jobs.title',
});
definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: 'JobsService.ListColleagues',
    redirect: { name: 'jobs-overview' },
});

const open = ref(false);
</script>

<template>
    <div class="min-h-full">
        <nav class="bg-base-700">
            <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="flex h-16 items-center justify-between">
                    <div class="flex items-center">
                        <div class="-ml-2 flex md:hidden">
                            <!-- Mobile menu button -->
                            <button
                                type="button"
                                class="relative inline-flex items-center justify-center rounded-md bg-base-500 p-2 text-base-200 hover:bg-base-400 hover:bg-opacity-75 hover:text-neutral focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2 focus:ring-offset-base-600"
                                @click="open = !open"
                            >
                                <span class="absolute -inset-0.5" />
                                <span class="sr-only">{{ $t('components.partials.sidebar.open_navigation') }}</span>
                                <MenuIcon v-if="!open" class="block h-5 w-5" aria-hidden="true" />
                                <CloseIcon v-else class="block h-5 w-5" aria-hidden="true" />
                            </button>
                        </div>
                        <div class="hidden md:block">
                            <div class="flex items-baseline space-x-2">
                                <template v-for="item in navigation" :key="item.name">
                                    <NuxtLink
                                        v-if="item.permission === undefined || can(item.permission)"
                                        :to="item.to"
                                        class="group flex shrink-0 flex-col items-center rounded-md p-3 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                                        active-class="bg-accent-100/20 text-neutral font-bold"
                                        exact-active-class="text-neutral"
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
                            class="group flex w-full flex-col items-center rounded-md p-2 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                            active-class="bg-accent-100/20 text-neutral font-bold"
                            exact-active-class="text-neutral"
                            aria-current-value="page"
                        >
                            {{ $t(item.name) }}
                        </NuxtLink>
                    </template>
                </div>
            </div>
        </nav>

        <main>
            <div class="mx-auto max-w-7xl py-4 sm:px-6 lg:px-8">
                <NuxtLayout name="blank">
                    <NuxtPage />
                </NuxtLayout>
            </div>
        </main>
    </div>
</template>

<script lang="ts" setup>
import {
    AccountGroupIcon,
    BriefcaseIcon,
    BulletinBoardIcon,
    CloseIcon,
    ListStatusIcon,
    MenuIcon,
    SchoolIcon,
    TimelineClockIcon,
} from 'mdi-vue3';
import type { DefineComponent } from 'vue';
import { type RoutesNamedLocations } from '@typed-router';
import type { Perms } from '~~/gen/ts/perms';

const navigation: { name: string; to: RoutesNamedLocations; permission?: Perms; icon: DefineComponent }[] = [
    {
        name: 'common.overview',
        to: { name: 'jobs-overview' },
        permission: 'JobsService.ListColleagues',
        icon: markRaw(BriefcaseIcon),
    },
    {
        name: 'pages.jobs.colleagues.title',
        to: { name: 'jobs-colleagues' },
        permission: 'JobsService.ListColleagues',
        icon: markRaw(AccountGroupIcon),
    },
    {
        name: 'common.activity',
        to: { name: 'jobs-activity' },
        permission: 'JobsService.GetColleague',
        icon: markRaw(BulletinBoardIcon),
    },
    {
        name: 'pages.jobs.timeclock.title',
        to: { name: 'jobs-timeclock' },
        permission: 'JobsTimeclockService.ListTimeclock',
        icon: markRaw(TimelineClockIcon),
    },
    {
        name: 'pages.jobs.qualifications.title',
        to: { name: 'jobs-qualifications' },
        permission: 'TODOService.TODOMethod',
        icon: markRaw(SchoolIcon),
    },
    {
        name: 'pages.jobs.conduct.title',
        to: { name: 'jobs-conduct' },
        permission: 'JobsConductService.ListConductEntries',
        icon: markRaw(ListStatusIcon),
    },
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
                    <div class="flex items-center md:overflow-x-auto">
                        <div class="-ml-2 flex md:hidden">
                            <!-- Mobile menu button -->
                            <button
                                type="button"
                                class="relative inline-flex items-center justify-center rounded-md bg-base-500 p-2 text-accent-200 hover:bg-base-400 hover:bg-opacity-75 hover:text-neutral focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2 focus:ring-offset-base-600"
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
                                    <span v-if="item.permission === undefined || can(item.permission)" class="flex-1">
                                        <NuxtLink
                                            v-slot="{ active }"
                                            :to="item.to"
                                            class="group flex shrink-0 items-center gap-2 rounded-md p-3 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                                            active-class="bg-accent-100/20 font-bold text-primary-300"
                                            aria-current-value="page"
                                            @click="open = false"
                                        >
                                            <component
                                                :is="item.icon"
                                                class="h-5 w-5"
                                                :class="active ? '' : 'group-hover:text-base-300'"
                                                aria-hidden="true"
                                            />
                                            {{ $t(item.name) }}
                                        </NuxtLink>
                                    </span>
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
                            v-slot="{ active }"
                            :to="item.to"
                            class="group flex shrink-0 items-center gap-2 rounded-md p-3 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                            active-class="bg-accent-100/20 font-bold text-primary-300"
                            aria-current-value="page"
                            @click="open = false"
                        >
                            <component
                                :is="item.icon"
                                :class="[active ? '' : 'group-hover:text-base-300', 'h-5 w-5']"
                                aria-hidden="true"
                            />
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

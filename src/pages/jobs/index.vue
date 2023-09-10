<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { CloseIcon, MenuIcon } from 'mdi-vue3';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';

const navigation: { name: string; to: RoutesNamedLocations; permission?: string }[] = [
    { name: 'Überblick', to: { name: 'jobs-index-overview' } },
    { name: 'Kollegen', to: { name: 'jobs-index-colleagues' } },
    { name: 'Anfragen', to: { name: 'jobs-index-requests' } },
    { name: 'Trainings', to: { name: 'jobs-index-trainings' } },
    { name: 'Stempeluhr', to: { name: 'jobs-index-timeclock' } },
    { name: 'Führungsregister', to: { name: 'jobs-index-conduct' }, permission: 'Jobs.ConductListEntries' },
];

useHead({
    title: 'pages.jobs.title',
});
definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: 'Jobs.View',
});

const route = useRoute();
</script>

<template>
    <div class="min-h-full">
        <Disclosure as="nav" class="bg-primary-600" v-slot="{ open }">
            <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="flex h-16 items-center justify-between">
                    <div class="flex items-center">
                        <div class="-ml-2 flex md:hidden">
                            <!-- Mobile menu button -->
                            <DisclosureButton
                                class="relative inline-flex items-center justify-center rounded-md bg-primary-600 p-2 text-primary-200 hover:bg-primary-500 hover:bg-opacity-75 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-primary-600"
                            >
                                <span class="absolute -inset-0.5" />
                                <span class="sr-only">{{ $t('components.partials.sidebar.open_menu') }}</span>
                                <MenuIcon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
                                <CloseIcon v-else class="block h-6 w-6" aria-hidden="true" />
                            </DisclosureButton>
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
                                        {{ item.name }}
                                    </NuxtLink>
                                </template>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <DisclosurePanel class="md:hidden">
                <div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
                    <template v-for="item in navigation" :key="item.name">
                        <DisclosureButton
                            v-if="item.permission === undefined || can(item.permission)"
                            as="NuxtLink"
                            :to="item.to"
                            class="text-white hover:bg-primary-500 hover:bg-opacity-75 block rounded-md px-3 py-2 text-base font-medium"
                            active-class="bg-primary-700 text-white"
                            aria-current-value="page"
                        >
                            {{ item.name }}
                        </DisclosureButton>
                    </template>
                </div>
            </DisclosurePanel>
        </Disclosure>

        <header class="bg-base-700 shadow">
            <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
                <h1 class="text-3xl font-bold leading-tight tracking-tight text-white">
                    {{ $t(route.meta.title ?? 'pages.jobs.title') }}
                </h1>
            </div>
        </header>
        <main>
            <div class="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
                <NuxtPage />
            </div>
        </main>
    </div>
</template>

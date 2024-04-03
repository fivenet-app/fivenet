<script lang="ts" setup>
import ExternalFileHolder from '~/components/partials/ExternalFileHolder.vue';

useHead({
    title: 'common.licenses',
});
definePageMeta({
    title: 'common.licenses',
    layout: 'landing',
    requiresAuth: false,
    showCookieOptions: true,
});

const licenses = [
    {
        label: 'FiveNet Project License',
        path: '/licenses/LICENSE',
    },
    {
        label: 'Frontend Licenses',
        path: '/licenses/frontend.txt',
    },
    {
        label: 'Sounds Licenses',
        path: '/licenses/sounds.txt',
    },
    {
        label: 'Backend Licenses',
        path: '/licenses/backend.txt',
    },
];
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <div class="flex h-full flex-col justify-between">
                <div>
                    <div class="hero relative isolate bg-primary-900 px-6 py-20 lg:px-8">
                        <div class="hero-overlay absolute left-0 top-0 z-[-1] size-full"></div>
                        <div class="mx-auto max-w-2xl text-center">
                            <h2 class="text-4xl font-bold tracking-tight sm:text-6xl">
                                {{ $t('common.licenses') }}
                            </h2>
                            <p class="mt-6 text-lg leading-8 text-gray-300">
                                {{ $t('components.about.licenses.subtitle') }}
                            </p>
                        </div>
                    </div>

                    <div>
                        <div class="mx-auto max-w-6xl px-2">
                            <div class="mx-auto max-w-7xl divide-y divide-neutral/10">
                                <dl class="mt-10 space-y-6 divide-y divide-neutral/10">
                                    <UAccordion :items="licenses" multiple>
                                        <template #item="{ item: license }">
                                            <ExternalFileHolder :path="license.path" />
                                        </template>
                                    </UAccordion>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>

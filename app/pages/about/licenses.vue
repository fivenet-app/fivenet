<script lang="ts" setup>
import ExternalFileHolder from '~/components/partials/ExternalFileHolder.vue';

useHead({
    title: 'common.licenses',
});

definePageMeta({
    title: 'common.licenses',
    layout: 'landing',
    requiresAuth: false,
    redirectIfAuthed: false,
    showCookieOptions: true,
});

const licenses = [
    {
        label: `FiveNet Project`,
        path: '/licenses/LICENSE.txt',
    },
    {
        label: `Frontend`,
        path: '/licenses/frontend.txt',
    },
    {
        label: `Backend`,
        path: '/licenses/backend.txt',
    },
    {
        label: `Sounds`,
        path: '/licenses/sounds.txt',
    },
];
</script>

<template>
    <UDashboardPanel>
        <div class="flex h-full flex-col justify-between">
            <div class="relative isolate px-6 py-20 lg:px-8">
                <div class="hero absolute inset-0 z-[-1] mask-[radial-gradient(100%_100%_at_top,white,transparent)]" />

                <div class="mx-auto max-w-2xl text-center">
                    <h2 class="text-4xl font-bold tracking-tight sm:text-6xl">
                        {{ $t('common.licenses') }}
                    </h2>
                    <p class="mt-6 text-lg leading-8">
                        {{ $t('pages.about.licenses.subtitle') }}
                    </p>
                </div>
            </div>

            <UContainer class="px-2">
                <div class="divide-neutral/10 mx-auto max-w-7xl divide-y">
                    <dl class="divide-neutral/10 mt-10 space-y-6 divide-y">
                        <UAccordion :items="licenses" type="multiple">
                            <template #content="{ item: license }">
                                <ExternalFileHolder :path="license.path" />
                            </template>
                        </UAccordion>
                    </dl>
                </div>
            </UContainer>
        </div>
    </UDashboardPanel>
</template>

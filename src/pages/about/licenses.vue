<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { ChevronDownIcon } from 'mdi-vue3';
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
        title: 'FiveNet Project License',
        path: '/licenses/LICENSE',
    },
    {
        title: 'Frontend Licenses',
        path: '/licenses/frontend.txt',
    },
    {
        title: 'Sounds Licenses',
        path: '/licenses/sounds.txt',
    },
    {
        title: 'Backend Licenses',
        path: '/licenses/backend.txt',
    },
];
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <div class="flex h-full flex-col justify-between">
                <div>
                    <div class="hero relative isolate bg-primary-900 px-6 py-24 sm:py-32 lg:px-8">
                        <div class="hero-overlay absolute left-0 top-0 z-[-1] size-full"></div>
                        <div class="mx-auto max-w-2xl text-center">
                            <h2 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                                {{ $t('common.licenses') }}
                            </h2>
                            <p class="mt-6 text-lg leading-8 text-gray-300">
                                {{ $t('components.about.licenses.subtitle') }}
                            </p>
                        </div>
                    </div>

                    <div class="bg-primary-900">
                        <div class="mx-auto max-w-7xl px-6 py-24 sm:py-32 lg:px-8 lg:py-20">
                            <div class="mx-auto max-w-7xl divide-y divide-neutral/10">
                                <dl class="mt-10 space-y-6 divide-y divide-neutral/10">
                                    <Disclosure
                                        v-for="license in licenses"
                                        :key="license.path"
                                        v-slot="{ open }"
                                        as="div"
                                        class="pt-6"
                                    >
                                        <dt>
                                            <DisclosureButton
                                                class="flex w-full items-start justify-between text-left text-neutral"
                                            >
                                                <span class="text-base font-semibold leading-7">{{ license.title }}</span>
                                                <span class="ml-6 flex h-7 items-center">
                                                    <ChevronDownIcon
                                                        :class="[open ? 'upsidedown' : '', 'size-5 transition-transform']"
                                                        aria-hidden="true"
                                                    />
                                                </span>
                                            </DisclosureButton>
                                        </dt>
                                        <DisclosurePanel as="dd" class="mt-2 pr-12">
                                            <p class="max-w-full">
                                                <ExternalFileHolder v-if="open" :path="license.path" />
                                            </p>
                                        </DisclosurePanel>
                                    </Disclosure>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>

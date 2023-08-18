<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { ChevronDownIcon } from 'mdi-vue3';
import Footer from '~/components/partials/Footer.vue';
import licensesBackend from '~/public/licenses-backend.txt?raw';
import licensesFrontend from '~/public/licenses-frontend.txt?raw';
import licenseFiveNet from '~~/LICENSE?raw';

useHead({
    title: 'common.licenses',
});
definePageMeta({
    title: 'common.licenses',
    requiresAuth: false,
});

const licenses = [
    {
        title: 'FiveNet License',
        content: licenseFiveNet,
    },
    {
        title: 'Frontend Licenses',
        content: licensesFrontend,
    },
    {
        title: 'Backend Licenses',
        content: licensesBackend,
    },
];
</script>

<template>
    <div class="h-full justify-between flex flex-col">
        <div>
            <div class="relative isolate bg-gray-900 px-6 py-24 sm:py-32 lg:px-8 hero">
                <div class="hero-overlay absolute z-[-1] h-full w-full top-0 left-0"></div>
                <div class="mx-auto max-w-2xl text-center">
                    <h2 class="text-4xl font-bold tracking-tight text-white sm:text-6xl">
                        {{ $t('common.licenses') }}
                    </h2>
                    <p class="mt-6 text-lg leading-8 text-gray-300">
                        {{ $t('pages.about.licenses.sub_title') }}
                    </p>
                </div>
            </div>

            <div class="bg-gray-900">
                <div class="mx-auto max-w-7xl px-6 py-24 sm:py-32 lg:px-8 lg:py-20">
                    <div class="mx-auto max-w-7xl divide-y divide-white/10">
                        <dl class="mt-10 space-y-6 divide-y divide-white/10">
                            <Disclosure v-for="license in licenses" as="div" class="pt-6" v-slot="{ open }">
                                <dt>
                                    <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                        <span class="text-base font-semibold leading-7">{{ license.title }}</span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                                aria-hidden="true"
                                            />
                                        </span>
                                    </DisclosureButton>
                                </dt>
                                <DisclosurePanel as="dd" class="mt-2 pr-12">
                                    <p class="max-w-full">
                                        <code
                                            class="max-w-full mt-2 p-4 whitespace-pre-line block bg-white text-black"
                                            v-text="license.content"
                                        ></code>
                                    </p>
                                </DisclosurePanel>
                            </Disclosure>
                        </dl>
                    </div>
                </div>
            </div>
        </div>
        <Footer />
    </div>
</template>

<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { ChevronDownIcon, LicenseIcon, LockIcon, MagnifyIcon, MapIcon } from 'mdi-vue3';
import '~/assets/css/herofull-pattern.css';
import PageFooter from '~/components/partials/PageFooter.vue';
import DiscordLogo from '~/components/partials/logos/DiscordLogo.vue';

const { t } = useI18n();

useHead({
    title: 'common.about',
});
definePageMeta({
    title: 'common.about',
    requiresAuth: false,
    showCookieOptions: true,
});

const discordLink = 'https://discord.gg/sWvkHuVQA5';
const repoLink = 'https://github.com/galexrt/fivenet';

const faqs = [
    {
        question: t('components.about.faq.one.question'),
        answer: t('components.about.faq.one.answer'),
    },
    {
        question: t('components.about.faq.two.question'),
        answer: t('components.about.faq.two.answer'),
    },
    {
        question: t('components.about.faq.three.question'),
        answer: t('components.about.faq.three.answer', { repoLink }),
    },
    {
        question: t('components.about.faq.four.question'),
        answer: t('components.about.faq.four.answer', { discordLink, repoLink }),
    },
] as { question: string; answer: string }[];
</script>

<template>
    <div class="flex h-full flex-col justify-between">
        <div>
            <div class="hero relative isolate bg-gray-900 px-6 py-12 sm:py-16 lg:px-8">
                <div class="hero-overlay absolute left-0 top-0 z-[-1] h-full w-full"></div>
                <div class="mx-auto max-w-2xl text-center">
                    <h2 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                        {{ $t('common.about') }}
                    </h2>
                    <p class="mt-6 text-lg leading-8 text-gray-300">
                        {{ $t('components.about.sub_title') }}
                    </p>
                </div>
            </div>

            <div class="relative isolate overflow-hidden px-6 py-12 sm:py-16 lg:overflow-visible lg:px-0">
                <div class="absolute inset-0 -z-10 overflow-hidden"></div>
                <div
                    class="mx-auto grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 lg:mx-0 lg:max-w-none lg:grid-cols-2 lg:items-start lg:gap-y-10"
                >
                    <div
                        class="lg:col-span-2 lg:col-start-1 lg:row-start-1 lg:mx-auto lg:grid lg:w-full lg:max-w-7xl lg:grid-cols-2 lg:gap-x-8 lg:px-8"
                    >
                        <div class="lg:pr-4">
                            <div class="lg:max-w-lg">
                                <p class="text-base font-semibold leading-7 text-primary-400">
                                    {{ $t('components.about.introduction.pre_title') }}
                                </p>
                                <h1 class="mt-2 text-3xl font-bold tracking-tight text-gray-400 sm:text-4xl">
                                    {{ $t('components.about.introduction.title') }}
                                </h1>
                                <p class="mt-6 text-xl leading-8 text-gray-100">
                                    {{ $t('components.about.introduction.content') }}
                                </p>
                            </div>
                        </div>
                    </div>
                    <div
                        class="-ml-12 -mt-12 p-12 lg:sticky lg:top-4 lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:overflow-hidden"
                    >
                        <img
                            class="w-[48rem] max-w-none rounded-xl bg-gray-900 shadow-xl ring-1 ring-gray-400/10 sm:w-[57rem]"
                            src="/images/app-screenshot.png"
                            alt="FiveNet Overview - Screenshot"
                        />
                    </div>
                    <div
                        class="lg:col-span-2 lg:col-start-1 lg:row-start-2 lg:mx-auto lg:grid lg:w-full lg:max-w-7xl lg:grid-cols-2 lg:gap-x-8 lg:px-8"
                    >
                        <div class="lg:pr-4">
                            <div class="max-w-xl text-base leading-7 text-gray-100 lg:max-w-lg">
                                <ul role="list" class="mt-8 space-y-8 text-gray-400">
                                    <li class="flex gap-x-3">
                                        <MagnifyIcon class="mt-1 h-5 w-5 flex-none text-primary-300" aria-hidden="true" />
                                        <span>
                                            <strong class="font-semibold text-gray-200">{{
                                                $t('components.about.introduction.feature_one.title')
                                            }}</strong>
                                            {{ $t('components.about.introduction.feature_one.content') }}
                                        </span>
                                    </li>
                                    <li class="flex gap-x-3">
                                        <LockIcon class="mt-1 h-5 w-5 flex-none text-primary-300" aria-hidden="true" />
                                        <span>
                                            <strong class="font-semibold text-gray-200">{{
                                                $t('components.about.introduction.feature_two.title')
                                            }}</strong>
                                            {{ $t('components.about.introduction.feature_two.content') }}
                                        </span>
                                    </li>
                                    <li class="flex gap-x-3">
                                        <MapIcon class="mt-1 h-5 w-5 flex-none text-primary-300" aria-hidden="true" />
                                        <span>
                                            <strong class="font-semibold text-gray-200">{{
                                                $t('components.about.introduction.feature_three.title')
                                            }}</strong>
                                            {{ $t('components.about.introduction.feature_three.content') }}
                                        </span>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="bg-gray-900">
                <div class="mx-auto max-w-7xl px-6 py-12 sm:py-16 lg:px-8 lg:py-10">
                    <div class="mx-auto max-w-4xl divide-y divide-neutral/10">
                        <h2 class="text-2xl font-bold leading-10 tracking-tight text-neutral">
                            {{ $t('components.about.faq.title') }}
                        </h2>
                        <dl class="mt-10 space-y-6 divide-y divide-neutral/10">
                            <Disclosure v-for="faq in faqs" :key="faq.question" v-slot="{ open }" as="div" class="pt-6">
                                <dt>
                                    <DisclosureButton class="flex w-full items-start justify-between text-left text-neutral">
                                        <span class="text-base font-semibold leading-7">{{ faq.question }}</span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-5 w-5 transition-transform']"
                                                aria-hidden="true"
                                            />
                                        </span>
                                    </DisclosureButton>
                                </dt>
                                <DisclosurePanel as="dd" class="mt-2 pr-12">
                                    <!-- eslint-disable-next-line vue/no-v-html -->
                                    <p class="text-base leading-7 text-gray-300" v-html="faq.answer"></p>
                                </DisclosurePanel>
                            </Disclosure>
                        </dl>
                    </div>
                </div>
            </div>
            <div class="relative bg-gray-900">
                <div class="mx-auto max-w-7xl px-6 py-12 sm:py-16 lg:px-8 lg:py-10">
                    <div class="mx-auto max-w-4xl">
                        <p class="mt-2 text-3xl font-bold tracking-tight text-neutral sm:text-4xl">
                            {{ $t('components.about.questions_or_issues.title') }}
                        </p>
                        <p class="mt-6 text-base leading-7 text-gray-300">
                            {{ $t('components.about.questions_or_issues.content') }}
                        </p>
                        <div class="mt-8">
                            <NuxtLink
                                :to="discordLink"
                                :external="true"
                                class="inline-flex items-center gap-x-2 rounded-md bg-neutral/10 px-3.5 py-2.5 text-sm font-semibold text-neutral shadow-sm hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-secondary-600"
                            >
                                <DiscordLogo class="-ml-0.5 h-5 w-5" aria-hidden="true" />
                                <span>
                                    {{ $t('components.about.join_discord') }}
                                </span>
                            </NuxtLink>
                        </div>
                    </div>
                </div>
            </div>
            <div class="relative bg-gray-900">
                <div class="mx-auto max-w-7xl px-6 py-12 sm:py-16 lg:px-8 lg:py-10">
                    <div class="mx-auto max-w-4xl">
                        <p class="mt-2 text-3xl font-bold tracking-tight text-neutral sm:text-4xl">
                            {{ $t('common.licenses') }}
                        </p>
                        <div class="mt-8">
                            <NuxtLink
                                :to="{ name: 'about-licenses' }"
                                class="inline-flex items-center gap-x-2 rounded-md bg-neutral/10 px-3.5 py-2.5 text-sm font-semibold text-neutral shadow-sm hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-secondary-600"
                            >
                                <LicenseIcon class="-ml-0.5 h-5 w-5" aria-hidden="true" />
                                <span>
                                    {{ $t('components.about.licenses_list') }}
                                </span>
                            </NuxtLink>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <PageFooter />
    </div>
</template>

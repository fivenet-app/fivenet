<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';

const { t } = useI18n();

useHead({
    title: 'pages.about.title',
});

definePageMeta({
    title: 'pages.about.title',
    layout: 'landing',
    requiresAuth: false,
    redirectIfAuthed: false,
    showCookieOptions: true,
});

const discordLink = 'https://discord.gg/ASRPPr8CeT';
const repoLink = 'https://github.com/fivenet-app/fivenet';

const faqs = [
    {
        label: t('pages.about.faq.one.question'),
        content: t('pages.about.faq.one.answer'),
    },
    {
        label: t('pages.about.faq.two.question'),
        content: t('pages.about.faq.two.answer'),
    },
    {
        label: t('pages.about.faq.three.question'),
        content: t('pages.about.faq.three.answer', { repoLink: repoLink }),
    },
    {
        label: t('pages.about.faq.four.question'),
        content: t('pages.about.faq.four.answer', { discordLink: discordLink, repoLink: repoLink }),
    },
] as { label: string; content: string }[];
</script>

<template>
    <UPage>
        <UDashboardPanel>
            <div class="flex flex-col justify-between">
                <div>
                    <div class="relative isolate px-6 py-20 lg:px-8">
                        <div
                            class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]"
                        />

                        <div class="mx-auto max-w-2xl text-center">
                            <h2 class="text-4xl font-bold tracking-tight sm:text-6xl">
                                {{ $t('pages.about.title') }}
                            </h2>
                            <p class="mt-6 text-lg leading-8">
                                {{ $t('pages.about.subtitle') }}
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
                                        <p class="text-primary-300 text-base font-semibold leading-7">
                                            {{ $t('pages.about.introduction.pre_title') }}
                                        </p>
                                        <h1 class="mt-2 text-3xl font-bold tracking-tight sm:text-4xl">
                                            {{ $t('pages.about.introduction.title') }}
                                        </h1>
                                        <p class="mt-6 text-xl leading-8 text-gray-100">
                                            {{ $t('pages.about.introduction.content') }}
                                        </p>
                                    </div>
                                </div>
                            </div>
                            <div
                                class="-ml-12 -mt-12 p-12 lg:sticky lg:top-4 lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:overflow-hidden"
                            >
                                <NuxtImg
                                    class="bg-background w-[48rem] max-w-none rounded-xl shadow-xl ring-1 ring-gray-400/10 sm:w-[57rem]"
                                    src="/images/screenshots/overview.png"
                                    alt="FiveNet Overview - Screenshot"
                                    loading="lazy"
                                />
                            </div>
                            <div
                                class="lg:col-span-2 lg:col-start-1 lg:row-start-2 lg:mx-auto lg:grid lg:w-full lg:max-w-7xl lg:grid-cols-2 lg:gap-x-8 lg:px-8"
                            >
                                <div class="lg:pr-4">
                                    <div class="max-w-xl text-base leading-7 text-gray-100 lg:max-w-lg">
                                        <ul class="mt-8 space-y-8" role="list">
                                            <li class="flex gap-x-3">
                                                <UIcon class="text-primary-300 mt-1 size-5 flex-none" name="i-mdi-magnify" />
                                                <span>
                                                    <span class="font-semibold text-gray-200">{{
                                                        $t('pages.about.introduction.feature_one.title')
                                                    }}</span>
                                                    {{ $t('pages.about.introduction.feature_one.content') }}
                                                </span>
                                            </li>
                                            <li class="flex gap-x-3">
                                                <UIcon class="text-primary-300 mt-1 size-5 flex-none" name="i-mdi-lock" />
                                                <span>
                                                    <span class="font-semibold text-gray-200">{{
                                                        $t('pages.about.introduction.feature_two.title')
                                                    }}</span>
                                                    {{ $t('pages.about.introduction.feature_two.content') }}
                                                </span>
                                            </li>
                                            <li class="flex gap-x-3">
                                                <UIcon class="text-primary-300 mt-1 size-5 flex-none" name="i-mdi-map" />
                                                <span>
                                                    <span class="font-semibold text-gray-200">{{
                                                        $t('pages.about.introduction.feature_three.title')
                                                    }}</span>
                                                    {{ $t('pages.about.introduction.feature_three.content') }}
                                                </span>
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="relative">
                        <div class="mx-auto max-w-7xl px-6 py-8">
                            <div class="mx-auto max-w-4xl">
                                <h2 class="mt-2 text-3xl font-bold tracking-tight sm:text-4xl">
                                    {{ $t('pages.about.faq.title') }}
                                </h2>
                                <dl class="mt-4">
                                    <UAccordion :items="faqs" multiple>
                                        <template #item="{ item: faq }">
                                            <UContainer>
                                                <!-- eslint-disable vue/no-v-html -->
                                                <p
                                                    class="text-base leading-7 text-gray-900 dark:text-white"
                                                    v-html="faq.content"
                                                ></p>
                                            </UContainer>
                                        </template>
                                    </UAccordion>
                                </dl>
                            </div>
                        </div>
                    </div>

                    <div class="relative">
                        <div class="mx-auto max-w-7xl px-6 py-8">
                            <div class="mx-auto max-w-4xl">
                                <h2 class="mt-2 text-3xl font-bold tracking-tight sm:text-4xl">
                                    {{ $t('pages.about.questions_or_issues.title') }}
                                </h2>
                                <p class="mt-6 text-base leading-7">
                                    {{ $t('pages.about.questions_or_issues.content') }}
                                </p>
                                <div class="mt-8">
                                    <UButton
                                        class="text-[#5865f2]"
                                        :to="discordLink"
                                        :external="true"
                                        variant="outline"
                                        block
                                        icon="i-simple-icons-discord"
                                    >
                                        {{ $t('pages.about.join_discord') }}
                                    </UButton>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="relative">
                        <div class="mx-auto max-w-7xl px-6 py-12 sm:py-16 lg:px-8 lg:py-10">
                            <div class="mx-auto max-w-4xl">
                                <p class="mt-2 text-3xl font-bold tracking-tight sm:text-4xl">
                                    {{ $t('common.licenses') }}
                                </p>
                                <div class="mt-8">
                                    <UButton variant="outline" block icon="i-mdi-license" :to="{ name: 'about-licenses' }">
                                        {{ $t('pages.about.licenses_list') }}
                                    </UButton>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </UDashboardPanel>
    </UPage>
</template>

<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';
import type { PageCardProps } from '@nuxt/ui';
import StarsBg from '~/components/landing/StarsBg.vue';

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

const introCards = computed<PageCardProps[]>(() => [
    {
        title: t('pages.about.introduction.feature_one.title'),
        description: t('pages.about.introduction.feature_one.content'),
        icon: 'i-mdi-magnify',
    },
    {
        title: t('pages.about.introduction.feature_two.title'),
        description: t('pages.about.introduction.feature_two.content'),
        icon: 'i-mdi-lock',
    },
    {
        title: t('pages.about.introduction.feature_three.title'),
        description: t('pages.about.introduction.feature_three.content'),
        icon: 'i-mdi-map',
    },
]);

const featureCards = computed<PageCardProps[]>(() => [
    {
        title: t('common.citizen', 2),
        description: t('pages.overview.features.citizens'),
        icon: 'i-mdi-account-multiple-outline',
        to: '/citizens',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        icon: 'i-mdi-car-outline',
        to: '/vehicles',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        icon: 'i-mdi-file-document-box-multiple-outline',
        to: '/documents',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        icon: 'i-mdi-briefcase-outline',
        to: '/jobs',
    },
    {
        title: t('common.calendar'),
        description: t('pages.overview.features.calendar'),
        icon: 'i-mdi-calendar-outline',
        to: '/calendar',
    },
    {
        title: t('common.mail', 2),
        description: t('pages.overview.features.mailer'),
        icon: 'i-mdi-email-outline',
        to: '/mail',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        icon: 'i-mdi-map-outline',
        to: '/livemap',
    },
    {
        title: t('common.dispatch_center'),
        description: t('pages.overview.features.centrum'),
        icon: 'i-mdi-car-emergency',
        to: '/dispatch',
    },
    {
        title: t('common.qualification', 2),
        description: t('pages.overview.features.qualifications'),
        icon: 'i-mdi-school-outline',
        to: '/qualifications',
    },
    {
        title: t('common.wiki'),
        description: t('pages.overview.features.wiki'),
        icon: 'i-mdi-book-open-variant-outline',
        to: '/wiki',
    },
]);

const faqs = computed(
    () =>
        [
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
                slot: 'question-3',
            },
            {
                label: t('pages.about.faq.four.question'),
                slot: 'question-4',
            },
        ] as { label: string; content?: string; slot?: string }[],
);
</script>

<template>
    <div class="relative min-h-dvh overflow-hidden">
        <div class="hero pointer-events-none absolute inset-0 z-[-1]" />

        <UPage :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
            <UPageHero
                :title="$t('pages.about.title')"
                :description="$t('pages.about.subtitle')"
                orientation="horizontal"
                :ui="{
                    container: 'py-24 sm:py-24 lg:py-24',
                    description: 'text-(--ui-text-highlighted)',
                    title: 'text-4xl sm:text-6xl',
                }"
            >
                <template #top>
                    <div
                        class="absolute left-1/2 size-60 -translate-x-1/2 -translate-y-80 transform rounded-full blur-[300px] sm:size-80 dark:bg-(--ui-primary)"
                    />
                </template>

                <template #default>
                    <NuxtImg
                        class="w-3xl max-w-none rounded-2xl"
                        src="/images/screenshots/overview.webp"
                        alt="FiveNet Overview - Screenshot"
                        loading="lazy"
                    />
                </template>
            </UPageHero>

            <div class="relative">
                <div
                    class="pointer-events-none inset-x-0 h-20"
                    :style="{ background: 'linear-gradient(to bottom, transparent, var(--ui-bg))' }"
                />

                <UPageSection
                    class="bg-default"
                    :title="$t('pages.about.introduction.title')"
                    :description="$t('pages.about.introduction.content')"
                    orientation="horizontal"
                    :ui="{ container: 'py-8 sm:py-8 lg:py-8' }"
                >
                    <template #headline>
                        <p class="text-sm font-semibold tracking-[0.24em] text-primary-300 uppercase">
                            {{ $t('pages.about.introduction.pre_title') }}
                        </p>
                    </template>

                    <div class="mt-8 flex flex-col gap-2">
                        <UPageCard v-for="item in introCards" :key="item.title?.toString()" v-bind="item" variant="subtle" />
                    </div>
                </UPageSection>

                <UPageSection
                    :title="$t('pages.about.features.title')"
                    :description="$t('pages.about.features.subtitle')"
                    class="bg-default !pt-0"
                    :ui="{ container: 'py-8 sm:py-8 lg:py-8' }"
                >
                    <UPageGrid class="gap-4 sm:grid-cols-2 xl:grid-cols-3">
                        <UPageCard v-for="item in featureCards" :key="item.title?.toString()" v-bind="item" />
                    </UPageGrid>
                </UPageSection>

                <UPageSection
                    :title="$t('pages.about.faq.title')"
                    class="bg-default"
                    :ui="{ container: 'py-8 sm:py-8 lg:py-8' }"
                >
                    <div class="mx-auto max-w-4xl">
                        <UAccordion :items="faqs" type="multiple" :ui="{ content: 'mb-2' }">
                            <template #content="{ item: faq }">
                                <UContainer>
                                    <!-- eslint-disable vue/no-v-html -->
                                    <p class="text-base leading-7 text-highlighted" v-html="faq.content"></p>
                                </UContainer>
                            </template>

                            <template #question-3>
                                <UContainer>
                                    <p class="text-base leading-7 text-highlighted">
                                        <NuxtLink class="underline" external :to="`${repoLink}/#readme`">
                                            {{ $t('pages.about.faq.three.click_here') }}
                                        </NuxtLink>
                                    </p>
                                </UContainer>
                            </template>

                            <template #question-4>
                                <UContainer>
                                    <p class="text-base leading-7 text-highlighted">
                                        <I18nT keypath="pages.about.faq.four.answer">
                                            <template #discordLink>
                                                <NuxtLink class="underline" external :to="discordLink">
                                                    {{ $t('pages.about.faq.four.discord_link') }}
                                                </NuxtLink>
                                            </template>

                                            <template #repoLink>
                                                <NuxtLink class="underline" external :to="repoLink">
                                                    {{ $t('pages.about.faq.four.repo_link') }}
                                                </NuxtLink>
                                            </template>
                                        </I18nT>
                                    </p>
                                </UContainer>
                            </template>
                        </UAccordion>
                    </div>
                </UPageSection>

                <UPageSection
                    :title="$t('common.license', 2)"
                    class="bg-default !pt-0"
                    :ui="{ container: 'py-8 sm:py-8 lg:py-8' }"
                >
                    <UPageCard variant="subtle" class="mx-auto max-w-4xl">
                        <p class="text-base leading-7 text-toned">
                            {{ $t('pages.about.licenses.subtitle') }}
                        </p>

                        <div class="mt-8">
                            <UButton icon="i-mdi-license" block to="/about/licenses" :label="$t('pages.about.licenses_list')" />
                        </div>
                    </UPageCard>
                </UPageSection>

                <USeparator />

                <UPageCTA
                    class="overflow-hidden"
                    :title="$t('pages.about.questions_or_issues.title')"
                    :description="$t('pages.about.questions_or_issues.content')"
                    :links="[
                        {
                            class: 'border-[#5865f2] bg-[#5865f2] text-white ring-[#5865f2]/50 hover:bg-[#5865f2]/10 focus:ring-[#5865f2]/50',
                            variant: 'outline',
                            size: 'xl',
                            icon: 'i-simple-icons-discord',
                            to: discordLink,
                            external: true,
                            label: $t('pages.about.join_discord'),
                            ui: { leadingIcon: 'size-10' },
                        },
                    ]"
                    variant="naked"
                >
                    <div
                        class="absolute left-1/2 size-40 -translate-x-1/2 -translate-y-80 transform rounded-full blur-[250px] sm:size-50 dark:bg-(--ui-primary)"
                    />

                    <StarsBg />
                </UPageCTA>
            </div>
        </UPage>
    </div>
</template>

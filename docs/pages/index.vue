<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';

const { data: page } = await useAsyncData('index', () => queryContent('/').findOne());

useHead({
    title: 'common.home',
});
definePageMeta({
    title: 'common.home',
    layout: 'landing',
    requiresAuth: false,
    showCookieOptions: true,
});

const { t } = useI18n();

const appVersion = __APP_VERSION__.split('-')[0];

const links = [{ label: t('common.docs'), icon: 'i-mdi-book-open-variant-outline', size: 'lg', to: '/getting-started' }];

const features = {
    title: t('docs.features.title'),
    links: [
        {
            label: t('docs.features.links.explore.label'),
            trailingIcon: 'i-mdi-arrow-right',
            color: 'gray',
            to: '/getting-started',
            size: 'lg',
        },
    ],
    items: [
        {
            title: t('docs.features.items.citizens.title'),
            description: t('docs.features.items.citizens.description'),
            icon: 'i-mdi-account-multiple-outline',
            to: '/user-guides/citizens',
        },
        {
            title: t('docs.features.items.vehicles.title'),
            description: t('docs.features.items.vehicles.description'),
            icon: 'i-mdi-car-outline',
            to: '/user-guides/vehicles',
        },
        {
            title: t('docs.features.items.documents.title'),
            description: t('docs.features.items.documents.description'),
            icon: 'i-mdi-file-document-box-multiple-outline',
            to: '/user-guides/documents',
        },
        {
            title: t('docs.features.items.jobs.title'),
            description: t('docs.features.items.jobs.description'),
            icon: 'i-mdi-briefcase-outline',
            to: '/user-guides/jobs',
        },
        {
            title: t('docs.features.items.livemap.title'),
            description: t('docs.features.items.livemap.description'),
            icon: 'i-mdi-map-outline',
            to: '/user-guides/livemap',
        },
        {
            title: t('docs.features.items.centrum.title'),
            description: t('docs.features.items.centrum.description'),
            icon: 'i-mdi-car-emergency',
            to: '/user-guides/centrum',
        },
        {
            title: t('docs.features.items.i18n.title'),
            description: t('docs.features.items.i18n.description'),
            icon: 'i-mdi-language',
        },
        {
            title: t('docs.features.items.nuxt3_ui.title'),
            description: t('docs.features.items.nuxt3_ui.description'),
            icon: 'i-simple-icons-nuxtdotjs',
            to: 'https://nuxt.com',
        },
    ],
};
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero :title="$t('pages.index.welcome')" :description="$t('pages.index.subtext')" :links="links">
                <template #headline>
                    <UButton
                        color="gray"
                        :to="`https://github.com/galexrt/fivenet/releases/tag/${appVersion}`"
                        :external="true"
                        :label="$t('pages.index.whats_new_in', { version: appVersion })"
                        trailing-icon="i-mdi-arrow-right"
                        size="xs"
                        class="rounded-full"
                    />
                </template>

                <ULandingLogos :title="$t('pages.index.logos')" align="center">
                    <ULink v-for="icon in page.logos.icons" :key="icon" variant="link">
                        <img
                            :src="`/images/communities/${icon.image}`"
                            :alt="icon.alt"
                            class="h-12 w-12 flex-shrink-0 text-gray-900 lg:h-20 lg:w-20 dark:text-white"
                        />
                    </ULink>
                </ULandingLogos>
            </ULandingHero>

            <ULandingSection :title="features.title" :links="features.links">
                <UPageGrid>
                    <ULandingCard v-for="(item, index) of features.items" :key="index" v-bind="item" />
                </UPageGrid>
            </ULandingSection>
        </div>
    </div>
</template>

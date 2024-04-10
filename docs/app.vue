<script setup lang="ts">
import type { NavItem, ParsedContent } from '@nuxt/content/dist/runtime/types';

const { seo } = useAppConfig();

const { data: navigation } = await useAsyncData<NavItem[]>('navigation', () => fetchContentNavigation(queryContent('/docs')), {
    default: () => [],
});
const { data: files } = useLazyFetch<ParsedContent[]>('/api/search.json', {
    default: () => [],
    server: false,
});

useHead({
    meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
    link: [{ rel: 'icon', href: '/favicon.ico' }],
    htmlAttrs: {
        lang: 'en',
    },
});

useSeoMeta({
    titleTemplate: `%s - ${seo?.siteName}`,
    ogSiteName: seo?.siteName,
    ogImage: 'https://docs-template.nuxt.dev/social-card.png',
    twitterImage: 'https://docs-template.nuxt.dev/social-card.png',
    twitterCard: 'summary_large_image',
});

provide('navigation', navigation);
</script>

<template>
    <div>
        <NuxtLoadingIndicator color="repeating-linear-gradient(to right, #55dde0 0%, #34cdfe 50%, #7161ef 100%)" />

        <UMain>
            <NuxtLayout>
                <NuxtPage />
            </NuxtLayout>
        </UMain>

        <Footer />

        <ClientOnly>
            <LazyUContentSearch :files="files" :navigation="navigation" />
        </ClientOnly>

        <UModals />
        <UNotifications />

        <CookieControl />
    </div>
</template>

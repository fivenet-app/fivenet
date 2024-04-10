<script setup lang="ts">
import type { NuxtError } from '#app';
import type { NavItem, ParsedContent } from '@nuxt/content/dist/runtime/types';

useSeoMeta({
    title: 'Page not found',
    description: 'We are sorry but this page could not be found.',
});

defineProps({
    error: {
        type: Object as PropType<NuxtError>,
        required: true,
    },
});

useHead({
    htmlAttrs: {
        lang: 'en',
    },
});

const { data: navigation } = await useAsyncData<NavItem[]>('navigation', () => fetchContentNavigation(), {
    default: () => [],
});

const { data: files } = useLazyFetch<ParsedContent[]>('/api/search.json', {
    default: () => [],
    server: false,
});

provide('navigation', navigation);
</script>

<template>
    <div>
        <LandingHeader />

        <UMain>
            <UContainer>
                <UPage>
                    <UPageError :error="error" :clear-button="{ label: $t('pages.notfound.go_back') }" />
                </UPage>
            </UContainer>
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

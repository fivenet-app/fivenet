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

const { data: navigation } = await useAsyncData<NavItem[]>('navigation', () => fetchContentNavigation(queryContent('/docs')), {
    default: () => [],
    transform: (items) => items[0]?.children ?? [],
});

console.log(navigation.value);

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
                    <UPageError :error="error" />
                </UPage>
            </UContainer>
        </UMain>

        <Footer />

        <ClientOnly>
            <LazyUContentSearch :files="files" :navigation="navigation" />
        </ClientOnly>

        <UNotifications />
    </div>
</template>

<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import * as locales from '@nuxt/ui/locale';
import NotificationProvider from '~/components/notifications/NotificationProvider.vue';
import CookieControl from '~/components/partials/CookieControl.vue';
import { useSettingsStore } from '~/stores/settings';

const { locale, t, finalizePendingLocaleChange } = useI18n();
const { initLocale } = useAppLocale();

useHead({
    htmlAttrs: {
        lang: locale,
    },
    titleTemplate: (title?: string) => (title ? `${title?.includes('.') ? t(title) : title} - FiveNet` : 'FiveNet'),
});

useSeoMeta({
    applicationName: 'FiveNet',
    title: 'FiveNet',
    ogTitle: 'FiveNet',
    ogImage: '/images/social-card.webp',
    twitterImage: '/images/social-card.webp',
    twitterCard: 'summary_large_image',
});

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const onBeforeEnter = async () => await finalizePendingLocaleChange();

const router = useRouter();
const route = router.currentRoute;

await initLocale();
</script>

<template>
    <UApp :locale="locales[locale]">
        <NuxtLoadingIndicator color="repeating-linear-gradient(to right, #55dde0 0%, #34cdfe 50%, #7161ef 100%)" />
        <NuxtRouteAnnouncer />

        <NuxtLayout>
            <NuxtPage :transition="{ onBeforeEnter }" />
        </NuxtLayout>

        <ClientOnly>
            <NotificationProvider />
        </ClientOnly>

        <CookieControl v-if="!nuiEnabled && route.meta.showCookieOptions !== undefined && route.meta.showCookieOptions" />
    </UApp>
</template>

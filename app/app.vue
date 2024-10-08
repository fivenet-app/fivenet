<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import CookieControl from '~/components/partials/CookieControl.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useSettingsStore } from '~/store/settings';
import { useAuthStore } from './store/auth';

const { t, setLocale, finalizePendingLocaleChange } = useI18n();

const appConfig = useAppConfig();

const toast = useToast();

const colorMode = useColorMode();

const color = computed(() => (colorMode.value === 'dark' ? '#111827' : 'white'));

useHead({
    htmlAttrs: {
        lang: 'en',
    },
    meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { key: 'theme-color', name: 'theme-color', content: color },
    ],
    link: [
        {
            rel: 'icon',
            type: 'image/png',
            href: '/favicon.png',
        },
    ],
    titleTemplate: (title?: string) => {
        if (title?.includes('.')) {
            title = t(title);
        }
        return title ? `${title} - FiveNet` : 'FiveNet';
    },
});

useSeoMeta({
    applicationName: 'FiveNet',
    title: 'FiveNet',
    ogTitle: 'FiveNet',
    ogImage: '/images/social-card.png',
    twitterImage: '/images/social-card.png',
    twitterCard: 'summary_large_image',
});

const settings = useSettingsStore();
const { locale: userLocale, isNUIAvailable, design, updateAvailable } = storeToRefs(settings);

if (APP_VERSION !== settings.version) {
    useLogger('⚙️ Settings').info('Resetting app data because new version has been detected', settings.version, APP_VERSION);

    useClipboardStore().$reset();
    useDocumentEditorStore().$reset();
    settings.setVersion(APP_VERSION);
}

// Set locale and theme colors in app config
appConfig.ui.primary = design.value.ui.primary;
appConfig.ui.gray = design.value.ui.gray;

if (userLocale.value !== null) {
    await setLocale(userLocale.value);
}

const onBeforeEnter = async () => {
    await finalizePendingLocaleChange();
};

onMounted(async () => {
    if (!import.meta.client) {
        return;
    }

    if (isNUIAvailable.value) {
        // NUI message handling
        window.addEventListener('message', onNUIMessage);
    }

    window.addEventListener('focusin', onFocusHandler, true);
    window.addEventListener('focusout', onFocusHandler, true);
});

onBeforeUnmount(async () => {
    if (!import.meta.client) {
        return;
    }

    if (isNUIAvailable.value) {
        // NUI message handling
        window.removeEventListener('message', onNUIMessage);
    }

    window.removeEventListener('focusin', onFocusHandler);
    window.removeEventListener('focusout', onFocusHandler);
});

watch(updateAvailable, async () => {
    if (!updateAvailable.value) {
        return;
    }

    toast.add({
        title: t('system.update_available.title', { version: updateAvailable.value }),
        description: t('system.update_available.content'),
        actions: [
            {
                label: t('common.refresh'),
                click: () => reloadNuxtApp({ persistState: false, force: true }),
            },
        ],
        icon: 'i-mdi-update',
        color: 'amber',
        timeout: 20000,
    });
});

const authStore = useAuthStore();
const { username } = storeToRefs(authStore);

// Use fivenet_authed cookie for basic browser-wide is logged in/out "signal"
const authedState = useCookie('fivenet_authed');
useIntervalFn(async () => refreshCookie('fivenet_authed'), 1750);

async function handleAuthedStateChange(): Promise<void> {
    if (!!authedState.value && username.value === null) {
        await authStore.chooseCharacter(undefined, true);
    } else if (!authedState.value && username.value !== null) {
        await navigateTo('/auth/logout');
    }
}

watch(authedState, handleAuthedStateChange);
handleAuthedStateChange();

const router = useRouter();
const route = router.currentRoute;
</script>

<template>
    <div>
        <NuxtLoadingIndicator color="repeating-linear-gradient(to right, #55dde0 0%, #34cdfe 50%, #7161ef 100%)" />

        <NuxtRouteAnnouncer />
        <NuxtLayout>
            <NuxtPage
                :transition="{
                    onBeforeEnter,
                }"
            />
        </NuxtLayout>

        <UNotifications />
        <UModals />
        <USlideovers />

        <ClientOnly>
            <LazyOverlaysSounds />
            <NotificationProvider />
        </ClientOnly>

        <CookieControl v-if="!isNUIAvailable && route.meta.showCookieOptions !== undefined && route.meta.showCookieOptions" />
    </div>
</template>

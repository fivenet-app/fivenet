<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import NotificationProvider from '~/components/notifications/NotificationProvider.vue';
import CookieControl from '~/components/partials/CookieControl.vue';
import { useAuthStore } from '~/stores/auth';
import { useClipboardStore } from '~/stores/clipboard';
import { useSettingsStore } from '~/stores/settings';
import BannerMessage from './components/partials/BannerMessage.vue';

const { locale, t, setLocale, finalizePendingLocaleChange } = useI18n();

const appConfig = useAppConfig();

const toast = useToast();

const colorMode = useColorMode();

const color = computed(() => (colorMode.value === 'dark' ? '#111827' : '#fff'));

const logger = useLogger('⚙️ Settings');

useHead({
    htmlAttrs: {
        lang: 'en',
    },
    meta: [{ key: 'theme-color', name: 'theme-color', content: color }],
    titleTemplate: (title?: string) => (title ? `${title?.includes('.') ? t(title) : title} - FiveNet` : 'FiveNet'),
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
const { getUserLocale, nuiEnabled, design, updateAvailable } = storeToRefs(settings);

if (APP_VERSION !== settings.version) {
    logger.info('Resetting app data because new version has been detected', settings.version, APP_VERSION);

    useClipboardStore().clear();
    useSearchesStore().clear();
    settings.setVersion(APP_VERSION);
}

// Set locale and theme colors in app config
async function setThemeColors(): Promise<void> {
    appConfig.ui.colors.primary = design.value.ui.primary;
    appConfig.ui.colors.neutral = design.value.ui.gray;
    setTabletColors(appConfig.ui.colors.primary, appConfig.ui.colors.neutral);
}
setThemeColors();
watch(design.value, setThemeColors);

async function setUserLocale(): Promise<void> {
    logger.info('Setting user locale to', getUserLocale.value);
    if (getUserLocale.value !== undefined) {
        locale.value = getUserLocale.value;
        await setLocale(getUserLocale.value);
    }
}
setUserLocale();
watch(getUserLocale, setUserLocale);

async function clickListener(event: MouseEvent): Promise<void> {
    if (!event.target || event.defaultPrevented) {
        return;
    }

    const element = event.target as HTMLElement;
    if (element.tagName.toLowerCase() !== 'a' && !element.hasAttribute('href')) {
        return;
    }
    const href = element.getAttribute('href');
    if (href?.startsWith('/') || href?.startsWith('#') || href === '') {
        return;
    }

    event.preventDefault();
    await navigateTo({
        name: 'dereferer',
        query: {
            target: href,
        },
    });
}

onMounted(async () => {
    if (!import.meta.client) {
        return;
    }

    if (nuiEnabled.value) {
        // NUI message handling
        window.addEventListener('message', onNUIMessage);
    }

    window.addEventListener('click', clickListener);
    window.addEventListener('focusin', onFocusHandler, true);
    window.addEventListener('focusout', onFocusHandler, true);
});

onBeforeUnmount(async () => {
    if (!import.meta.client) {
        return;
    }

    if (nuiEnabled.value) {
        // NUI message handling
        window.removeEventListener('message', onNUIMessage);
    }

    window.removeEventListener('click', clickListener);
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
                onClick: () => reloadNuxtApp({ persistState: false, force: true }),
            },
        ],
        icon: 'i-mdi-update',
        color: 'primary',
        duration: 20000,
        close: {
            disabled: true,
        },
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

const onBeforeEnter = async () => {
    await finalizePendingLocaleChange();
};

const router = useRouter();
const route = router.currentRoute;
</script>

<template>
    <UApp>
        <NuxtLoadingIndicator color="repeating-linear-gradient(to right, #55dde0 0%, #34cdfe 50%, #7161ef 100%)" />
        <NuxtRouteAnnouncer />

        <NuxtLayout>
            <NuxtPage :transition="{ onBeforeEnter }" />
        </NuxtLayout>

        <BannerMessage
            v-if="appConfig.system.bannerMessageEnabled && appConfig.system.bannerMessage"
            :message="appConfig.system.bannerMessage"
        />

        <ClientOnly>
            <NotificationProvider />
        </ClientOnly>

        <CookieControl v-if="!nuiEnabled && route.meta.showCookieOptions !== undefined && route.meta.showCookieOptions" />
    </UApp>
</template>

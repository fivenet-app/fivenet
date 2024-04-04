<script lang="ts" setup>
import { localize, setLocale as veeValidateSetLocale } from '@vee-validate/i18n';
import de from '@vee-validate/i18n/dist/locale/de.json';
import en from '@vee-validate/i18n/dist/locale/en.json';
import { configure } from 'vee-validate';
import { useClipboardStore } from '~/store/clipboard';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useSettingsStore } from '~/store/settings';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';

const { t, locale, finalizePendingLocaleChange } = useI18n();

const settings = useSettingsStore();
const { isNUIAvailable, updateAvailable, design } = storeToRefs(settings);

const route = useRoute();

const { locale: cookieLocale } = useCookieControl();

const colorMode = useColorMode();

const color = computed(() => (colorMode.value === 'dark' ? '#111827' : 'white'));

useHead({
    meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { key: 'theme-color', name: 'theme-color', content: color },
    ],
    htmlAttrs: {
        class: 'h-full scrollbar-thin scrollbar-thumb-sky-700 scrollbar-track-sky-300',
        lang: 'en',
    },
    bodyAttrs: {
        class: 'h-full',
    },
    titleTemplate: (title?: string) => {
        if (title?.includes('.')) {
            title = t(title);
        }
        return title ? `${title} - FiveNet` : 'FiveNet';
    },
});
useSeoMeta({
    applicationName: 'FiveNet',
    ogImage: '/images/social-card.png',
    twitterImage: '/images/social-card.png',
});

if (__APP_VERSION__ !== settings.version) {
    console.info('Resetting app data because new version has been detected', settings.version, __APP_VERSION__);
    useClipboardStore().$reset();
    useDocumentEditorStore().$reset();
    settings.setVersion(__APP_VERSION__);
}

configure({
    generateMessage: localize({
        en,
        de,
    }),
});

// Set user setting locale on load of app
if (settings.locale !== null) {
    locale.value = settings.locale;
} else {
    settings.locale = locale.value;
}
setLocaleGlobally(locale.value);

async function setLocaleGlobally(locale: string): Promise<void> {
    settings.setLocale(locale);
    veeValidateSetLocale(locale);

    // Cookie Banner Locale handling
    switch (locale.split('-', 1)[0]) {
        case 'de':
            cookieLocale.value = 'de';
            break;
        default:
            cookieLocale.value = 'en';
            break;
    }
}

async function onBeforeEnter(): Promise<void> {
    await finalizePendingLocaleChange();
}

watch(locale, () => setLocaleGlobally(locale.value));

// NUI message handling
if (isNUIAvailable.value) {
    onMounted(async () => window.addEventListener('message', onNUIMessage));
    onBeforeUnmount(async () => window.removeEventListener('message', onNUIMessage));
}

const toast = useToast();

watch(updateAvailable, () => {
    if (!updateAvailable.value) {
        return;
    }

    toast.add({
        title: t('system.update_available.title', { version: updateAvailable }),
        description: t('system.update_available.content'),
        actions: [
            {
                label: t('common.refresh'),
                click: () => reloadNuxtApp({ persistState: false, force: true }),
            },
        ],
        icon: 'i-mdi-update',
        color: 'amber',
    });
});
</script>

<template>
    <div>
        <NuxtLoadingIndicator />

        <NuxtLayout>
            <NuxtPage
                :transition="{
                    name: 'page',
                    mode: 'out-in',
                    onBeforeEnter,
                }"
            />
        </NuxtLayout>

        <UNotifications />
        <UModals />
        <USlideovers />

        <NotificationProvider />

        <CookieControl
            v-if="!isNUIAvailable && route.meta.showCookieOptions !== undefined && route.meta.showCookieOptions"
            :locale="cookieLocale"
        />
    </div>
</template>

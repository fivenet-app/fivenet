<script lang="ts" setup>
import { localize, setLocale as veeValidateSetLocale } from '@vee-validate/i18n';
import de from '@vee-validate/i18n/dist/locale/de.json';
import en from '@vee-validate/i18n/dist/locale/en.json';
import { UpdateIcon } from 'mdi-vue3';
import { configure } from 'vee-validate';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useConfigStore } from '~/store/config';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { JOB_THEME_KEY, useSettingsStore } from '~/store/settings';
import { useAuthStore } from '~/store/auth';

const { t, locale, finalizePendingLocaleChange } = useI18n();

const authStore = useAuthStore();
const { jobProps } = storeToRefs(authStore);

const configStore = useConfigStore();
const { isNUIAvailable, updateAvailable } = storeToRefs(configStore);

const settings = useSettingsStore();
const { theme } = storeToRefs(settings);

const route = useRoute();

const { locale: cookieLocale } = useCookieControl();

useHead({
    htmlAttrs: {
        class: () =>
            (theme.value === JOB_THEME_KEY ? jobProps.value?.theme ?? 'defaultTheme' : theme.value) +
            ' scrollbar-thin scrollbar-thumb-sky-700 scrollbar-track-sky-300',
        lang: 'en',
    },
    bodyAttrs: {
        class: 'bg-body-color h-full overflow-y-hidden',
    },
    titleTemplate: (title?: string) => {
        if (title?.includes('.')) {
            title = t(title);
        }
        return title ? `${title} - FiveNet` : 'FiveNet';
    },
});
useSeoMeta({
    ogImage: '/images/open-graph-image.png',
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
locale.value = settings.locale;
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

// Open update available confirm dialog
const open = ref(false);

watch(updateAvailable, () => (open.value = true));
/* eslint vue/no-multiple-template-root: "off" */
</script>

<template>
    <NuxtLayout>
        <NuxtPage
            :transition="{
                name: 'page',
                mode: 'out-in',
                onBeforeEnter,
            }"
        />
    </NuxtLayout>

    <CookieControl
        v-if="!isNUIAvailable && route.meta.showCookieOptions !== undefined && route.meta.showCookieOptions"
        :locale="cookieLocale"
    />

    <ConfirmDialog
        v-if="updateAvailable !== false"
        :open="open"
        :cancel="() => (open = false)"
        :confirm="() => reloadNuxtApp({ persistState: false, force: true })"
        :title="$t('system.update_available.title', { version: updateAvailable })"
        :description="$t('system.update_available.content')"
        :icon="markRaw(UpdateIcon)"
        @close="open = false"
    />
</template>

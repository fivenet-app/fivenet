<script lang="ts" setup>
import { Locale } from '@dargmuesli/nuxt-cookie-control/dist/runtime/types';
import { localize, setLocale as veeValidateSetLocale } from '@vee-validate/i18n';
import de from '@vee-validate/i18n/dist/locale/de.json';
import en from '@vee-validate/i18n/dist/locale/en.json';
import { UpdateIcon } from 'mdi-vue3';
import { configure } from 'vee-validate';
import { useClipboardStore } from '~/store/clipboard';
import { useConfigStore } from '~/store/config';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useNotificationsStore } from '~/store/notifications';
import { useSettingsStore } from '~/store/settings';
import ConfirmDialog from './components/partials/ConfirmDialog.vue';

const { t, setLocale } = useI18n();

const configStore = useConfigStore();
const { loadConfig } = configStore;
const { clientConfig, updateAvailable } = storeToRefs(configStore);

const route = useRoute();

useHead({
    htmlAttrs: {
        class: 'h-full',
        lang: 'en',
    },
    bodyAttrs: {
        class: 'bg-body-color h-full overflow-hidden',
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

await loadConfig();

const userSettings = useSettingsStore();
if (__APP_VERSION__ != userSettings.version) {
    console.info('Resetting app data because new version has been detected', userSettings.version, __APP_VERSION__);
    useClipboardStore().$reset();
    useDocumentEditorStore().$reset();
    useNotificationsStore().$reset();
    userSettings.setVersion(__APP_VERSION__);
}

// Set user setting locale on load of app
setLocale(userSettings.locale);

configure({
    generateMessage: localize({
        en,
        de,
    }),
});
veeValidateSetLocale(userSettings.locale);

// Cookie Banner Locale handling
const cookieLocale = ref<Locale>('en');
switch (userSettings.locale.split('-', 1)[0]) {
    case 'de':
        cookieLocale.value = 'de';
        break;
    default:
        cookieLocale.value = 'en';
        break;
}

const open = ref(false);

function triggerUpdate(): void {
    location.reload();
}

// Open update available confirm dialog
watch(updateAvailable, () => (open.value = true));
</script>

<template>
    <NuxtLayout>
        <NuxtPage
            :transition="{
                name: 'page',
                mode: 'out-in',
            }"
        />
    </NuxtLayout>
    <CookieControl
        v-if="!clientConfig.NUIEnabled && route.meta.showCookieOptions !== undefined && route.meta.showCookieOptions"
        :locale="cookieLocale"
    />

    <ConfirmDialog
        v-if="updateAvailable !== false"
        :open="open"
        @close="open = false"
        :cancel="() => (open = false)"
        :confirm="() => triggerUpdate()"
        :title="$t('system.update_available.title', [updateAvailable])"
        :description="$t('system.update_available.content')"
        :icon="UpdateIcon"
    />
</template>

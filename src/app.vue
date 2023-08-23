<script lang="ts" setup>
import { Locale } from '@dargmuesli/nuxt-cookie-control/dist/runtime/types';
import { localize, setLocale as veeValidateSetLocale } from '@vee-validate/i18n';
import de from '@vee-validate/i18n/dist/locale/de.json';
import en from '@vee-validate/i18n/dist/locale/en.json';
import { configure } from 'vee-validate';
import { useClipboardStore } from '~/store/clipboard';
import { useConfigStore } from '~/store/config';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useNotificationsStore } from '~/store/notifications';
import { useSettingsStore } from '~/store/settings';

const { t, setLocale } = useI18n();

const configStore = useConfigStore();
const { loadConfig } = configStore;
const { clientConfig } = storeToRefs(configStore);

const route = useRoute();

useHead({
    htmlAttrs: {
        class: 'h-full bg-base-900',
        lang: 'en',
    },
    bodyAttrs: {
        class: 'h-full overflow-hidden',
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

loadConfig();

const userSettings = useSettingsStore();
if (__APP_VERSION__ != userSettings.getVersion) {
    console.info('Resetting app data because new version detected', userSettings.getVersion, __APP_VERSION__);
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

if (route.query?.nui) {
    clientConfig.value.NUIEnabled = true;
    clientConfig.value.NUIResourceName = route.query?.nui as string;
    console.info('Enabled NUI integration! Resource Name:', clientConfig.value.NUIResourceName);
}
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
        :locale="cookieLocale"
        v-if="useRoute().meta.showCookieOptions !== undefined && useRoute().meta.showCookieOptions"
    />
</template>

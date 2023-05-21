<script lang="ts" setup>
import { loadConfig } from '~/config';
import { useUserSettingsStore } from './store/usersettings';
import { NuxtError } from 'nuxt/app';
import { configure } from 'vee-validate';
import { localize, setLocale as veeValidateSetLocale } from '@vee-validate/i18n';
import en from '@vee-validate/i18n/dist/locale/en.json';
import de from '@vee-validate/i18n/dist/locale/de.json';

const { t, setLocale } = useI18n();

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

try {
    await loadConfig();
} catch (e) {
    showError(e as NuxtError);
}

const userSettings = useUserSettingsStore();

// Set user setting locale on load of app
setLocale(userSettings.locale);

configure({
  generateMessage: localize({
    en,
    de,
  }),
});

veeValidateSetLocale(userSettings.locale);
</script>

<template>
    <NuxtLayout>
        <NuxtPage :transition="{
            name: 'page',
            mode: 'out-in'
        }" />
    </NuxtLayout>
</template>

<script lang="ts" setup>
import { loadConfig } from './config';
import { useUserSettingsStore } from './store/usersettings';

const { t, setLocale } = useI18n();

useHead({
    htmlAttrs: {
        class: 'h-full bg-base-900',
        lang: 'en',
    },
    bodyAttrs: {
        class: 'h-full overflow-hidden',
    },
    titleTemplate: (title) => {
        if (title?.includes('.')) {
            title = t(title);
        }
        return title ? `${title} - FiveNet` : 'FiveNet';
    },
});

await loadConfig();

const store = useUserSettingsStore();

// Set locale on load
const locale = computed(() => store.$state.locale);
setLocale(locale.value);
</script>

<template>
    <NuxtLayout>
        <NuxtPage :transition="{
            name: 'page',
            mode: 'out-in'
        }" />
    </NuxtLayout>
</template>

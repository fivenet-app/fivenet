<script lang="ts" setup>
import { loadConfig } from './config';
import { useUserSettingsStore } from './store/usersettings';

useHead({
    htmlAttrs: {
        class: 'h-full bg-base-900',
        lang: 'en',
    },
    bodyAttrs: {
        class: 'h-full overflow-hidden',
    },
    titleTemplate: (titleChunk) => {
        return titleChunk ? `${titleChunk} - FiveNet` : 'FiveNet';
    },
});

await loadConfig();

const store = useUserSettingsStore();
const { setLocale } = useI18n();

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

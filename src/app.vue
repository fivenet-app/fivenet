<script lang="ts" setup>
import { loadConfig } from './config';

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
</script>

<template>
    <NuxtLayout>
        <NuxtPage :transition="{
            name: 'page',
                mode: 'out-in'
        }" />
    </NuxtLayout>
</template>

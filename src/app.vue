<script lang="ts" setup>
import { loadConfig } from './config';
import { NotificationProvider } from '~/components/notification';
import Sidebar from '~/components/partials/Sidebar.vue';

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
    <NotificationProvider>
        <Sidebar>
            <NuxtLayout>
                <NuxtPage :transition="{
                    name: 'page',
                    mode: 'out-in'
                }" />
            </NuxtLayout>
        </Sidebar>
    </NotificationProvider>
</template>

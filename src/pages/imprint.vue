<script lang="ts" setup>
import { useTimeoutFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
import HeroPage from '~/components/partials/HeroPage.vue';

useHead({
    title: 'common.imprint',
});
definePageMeta({
    title: 'common.imprint',
    requiresAuth: false,
    showCookieOptions: true,
});

const { links } = useAppConfig();

if (links.imprint === undefined) {
    navigateTo({ name: 'index' });
} else {
    useTimeoutFn(() => navigateTo(links.imprint!, { external: true }), 1750);
}
</script>

<template>
    <HeroPage>
        <template #default>
            <h1 class="text-5xl font-bold tracking-tight text-neutral sm:text-6xl">
                {{ $t('common.redirecting_to', [$t('common.imprint')]) }}
            </h1>
            <h2 class="mt-4 inline-flex items-center text-4xl text-neutral">
                <LoadingIcon class="w-20 h-20 animate-spin" aria-hidden="true" />
            </h2>
        </template>
    </HeroPage>
</template>

<script lang="ts" setup>
import { useTimeoutFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import ContentHeroFull from '~/components/partials/ContentHeroFull.vue';
import PageFooter from '~/components/partials/PageFooter.vue';

useHead({
    title: 'common.imprint',
});
definePageMeta({
    title: 'common.imprint',
    requiresAuth: false,
    showCookieOptions: true,
});

const { links } = useAppConfig();

useTimeoutFn(() => navigateTo(links.imprint, { external: true }), 2000);
</script>

<template>
    <div class="h-full flex flex-col">
        <ContentHeroFull>
            <ContentCenterWrapper class="mx-auto max-w-2xl text-center">
                <div class="px-5 sm:px-0">
                    <h1 class="text-5xl font-bold tracking-tight text-neutral sm:text-6xl">
                        {{ $t('common.redirecting_to', [$t('common.imprint')]) }}
                    </h1>
                    <h2 class="mt-4 inline-flex items-center text-4xl text-neutral">
                        <LoadingIcon class="w-20 h-20 animate-spin" aria-hidden="true" />
                    </h2>
                </div>
            </ContentCenterWrapper>
        </ContentHeroFull>

        <PageFooter />
    </div>
</template>

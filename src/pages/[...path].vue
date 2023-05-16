<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';

useHead({
    title: 'pages.notfound.title',
});
definePageMeta({
    title: 'pages.notfound.title',
    requiresAuth: false,
});

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);
</script>

<template>
    <HeroFull>
        <ContentCenterWrapper class="max-w-2xl mx-auto text-center">
            <h1 class="text-5xl font-bold text-neutral">
                {{ $t('pages.notfound.page_not_found') }}
            </h1>
            <h2 class="text-4xl text-neutral">
                {{ $t('pages.notfound.error') }}
            </h2>
            <p class="py-6 text-neutral">
                {{ $t('pages.notfound.fun_error') }}
            </p>

            <NuxtLink v-if="accessToken" :to="{ name: 'overview' }"
                class="rounded-md w-60 bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                {{ $t('common.overview') }}
            </NuxtLink>
            <NuxtLink v-else :to="{ name: 'auth-login' }"
                class="rounded-md w-60 bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                {{ $t('common.login') }}
            </NuxtLink>
        </ContentCenterWrapper>
    </HeroFull>
</template>

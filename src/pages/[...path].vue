<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'pages.notfound.title',
});
definePageMeta({
    title: 'pages.notfound.title',
    requiresAuth: false,
    showCookieOptions: true,
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

            <NuxtLink
                v-if="accessToken"
                :to="{ name: 'overview' }"
                class="rounded-md w-60 bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
            >
                {{ $t('common.overview') }}
            </NuxtLink>
            <NuxtLink
                v-else
                :to="{ name: 'auth-login' }"
                class="w-48 max-w-96 rounded-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
            >
                {{ $t('pages.auth.login.menu_item') }}
            </NuxtLink>
        </ContentCenterWrapper>
    </HeroFull>
</template>

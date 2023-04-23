<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import HeroFull from '~/components/partials/HeroFull.vue';
import { computed } from 'vue';

useHead({
    title: '404 Not Found',
});
definePageMeta({
    title: '404 Not Found',
    requiresAuth: false,
});

const store = useAuthStore();
const accessToken = computed(() => store.$state.accessToken);
</script>

<template>
    <HeroFull>
        <h1 class="text-5xl font-bold text-neutral">{{ $t('pages.404.page_not_found') }}</h1>
        <h2 class="text-4xl text-neutral">{{ $t('pages.404.error') }}</h2>
        <p class="py-6 text-neutral">{{ $t('pages.404.fun_error') }}</p>
        <NuxtLink v-if="accessToken" :to="{ name: 'overview' }" class="btn btn-primary">{{ $t('common.overview') }}
        </NuxtLink>
        <NuxtLink v-else :to="{ name: 'auth-login' }" class="btn btn-primary">${{ $t('common.login') }}</NuxtLink>
    </HeroFull>
</template>

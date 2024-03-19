<script lang="ts" setup>
import { useTimeoutFn } from '@vueuse/core';
import { useAuthStore } from '~/store/auth';
import HeroPage from '~/components/partials/HeroPage.vue';

useHead({
    title: 'common.logout',
});
definePageMeta({
    title: 'common.logout',
    requiresAuth: true,
    authOnlyToken: true,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { doLogout } = authStore;

function redirect(): void {
    useTimeoutFn(async () => {
        const route = useRoute();
        if (route.name === 'auth-logout') {
            await navigateTo({ name: 'index' });
        }
    }, 1500);
}

onMounted(async () => {
    try {
        await doLogout();
    } finally {
        redirect();
    }
});
</script>

<template>
    <HeroPage>
        <template #default>
            <h2 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                {{ $t('components.auth.logout.header') }}
            </h2>
            <p class="mt-6 text-lg leading-8 text-gray-300">
                {{ $t('components.auth.logout.subtitle') }}
            </p>
        </template>
    </HeroPage>
</template>

<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

useHead({
    title: 'common.logout',
});
definePageMeta({
    title: 'common.logout',
    layout: 'auth',
    requiresAuth: true,
    authOnlyToken: true,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { doLogout } = authStore;

onMounted(async () => {
    try {
        await doLogout();
    } finally {
        useTimeoutFn(async () => {
            const route = useRoute();
            if (route.name === 'auth-logout') {
                await navigateTo({ name: 'index' });
            }
        }, 1500);
    }
});
</script>

<template>
    <UPage class="w-full flex-1">
        <ULandingHero :title="$t('pages.auth.logout.header')" :description="$t('pages.auth.logout.subtitle')">
            <template #headline>
                <FiveNetLogo class="mx-auto h-36 w-auto" />
            </template>
        </ULandingHero>
    </UPage>
</template>

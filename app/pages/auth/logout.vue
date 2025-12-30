<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { useAuthStore } from '~/stores/auth';

useHead({
    title: 'common.logout',
});

definePageMeta({
    title: 'common.logout',
    layout: 'auth',
    requiresAuth: true,
    authTokenOnly: true,
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
                if (route.query.redirect) {
                    const redirect = route.query.redirect as string;
                    if (redirect !== '/') {
                        // @ts-expect-error 404 handler will handle wrong URLs..
                        await navigateTo(redirect);
                        return;
                    }
                }

                await navigateTo('/');
            }
        }, 1500);
    }
});
</script>

<template>
    <UPage class="w-full flex-1">
        <UPageHero :title="$t('pages.auth.logout.header')" :description="$t('pages.auth.logout.subtitle')">
            <template #headline>
                <FiveNetLogo class="mx-auto h-auto w-20" />
            </template>
        </UPageHero>
    </UPage>
</template>

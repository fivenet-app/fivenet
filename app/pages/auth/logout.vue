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
    <UPageCard
        class="w-full max-w-md shrink-0 bg-white/75 backdrop-blur-sm dark:bg-white/5"
        :description="$t('pages.auth.logout.subtitle')"
        :ui="{ header: 'w-full' }"
    >
        <template #header>
            <div class="w-full space-y-2">
                <FiveNetLogo class="mx-auto h-auto w-20" />

                <h2 class="text-center text-3xl">
                    {{ $t('pages.auth.logout.header') }}
                </h2>
            </div>
        </template>
    </UPageCard>
</template>

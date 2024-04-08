<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import PageFooter from '~/components/partials/PageFooter.vue';
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
                navigateTo({ name: 'index' });
            }
        }, 1500);
    }
});
</script>

<template>
    <div class="flex size-full flex-col">
        <div class="hero w-full flex-1">
            <UContainer class="h-full bg-black/50">
                <UPage>
                    <ULandingHero :title="$t('pages.auth.logout.header')" :description="$t('components.auth.logout.subtitle')">
                        <template #headline>
                            <FiveNetLogo class="mx-auto h-36 w-auto" />
                        </template>
                    </ULandingHero>
                </UPage>
            </UContainer>
        </div>

        <PageFooter />
    </div>
</template>

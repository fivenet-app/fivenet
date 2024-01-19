<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import PageFooter from '~/components/partials/PageFooter.vue';
import ContentHeroFull from '~/components/partials/ContentHeroFull.vue';
import { useAuthStore } from '~/store/auth';

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
    setTimeout(async () => {
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
    <div class="flex h-full flex-col justify-between">
        <ContentHeroFull>
            <ContentCenterWrapper class="mx-auto max-w-2xl text-center">
                <h2 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                    {{ $t('components.auth.logout.header') }}
                </h2>
                <p class="mt-6 text-lg leading-8 text-gray-300">
                    {{ $t('components.auth.logout.subtitle') }}
                </p>
            </ContentCenterWrapper>
        </ContentHeroFull>
        <PageFooter />
    </div>
</template>

<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import Footer from '~/components/partials/Footer.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';

useHead({
    title: 'common.logout',
});
definePageMeta({
    title: 'common.logout',
    requiresAuth: true,
    authOnlyToken: true,
});

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const notifications = useNotificationsStore();

const { accessToken } = storeToRefs(authStore);
const { clearAuthInfo } = authStore;

const { t } = useI18n();

function redirect(): void {
    setTimeout(async () => {
        const route = useRoute();
        if (route.name === 'auth-logout') {
            await navigateTo({ name: 'index' });
        }
    }, 1500);
}

onBeforeMount(async () => {
    clearAuthInfo();

    if (!accessToken.value) {
        redirect();
        return;
    }

    try {
        const call = $grpc.getAuthClient().logout({});
        const { status } = await call;

        if (await $grpc.handleError(status)) {
            throw new Error(status.detail);
        }
    } catch (e) {
        notifications.dispatchNotification({
            title: t('notifications.auth.error_logout.title'),
            content: t('notifications.auth.error_logout.content', [e]),
            type: 'error',
        });
    }

    redirect();
});
</script>

<template>
    <div class="h-full justify-between flex flex-col">
        <HeroFull>
            <ContentCenterWrapper class="max-w-2xl mx-auto text-center">
                <h2 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                    {{ $t('pages.auth.logout.header') }}
                </h2>
                <p class="mt-6 text-lg leading-8 text-gray-300">
                    {{ $t('pages.auth.logout.subtitle') }}
                </p>
            </ContentCenterWrapper>
        </HeroFull>
        <Footer />
    </div>
</template>

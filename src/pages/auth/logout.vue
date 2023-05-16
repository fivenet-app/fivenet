<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import { LogoutRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import Footer from '~/components/partials/Footer.vue';

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
        if (route.name == 'auth-logout') {
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
        await $grpc.getAuthClient()
            .logout(new LogoutRequest(), null);
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);

        const err = e as RpcError;
        notifications.dispatchNotification({
            title: t('notifications.auth.error_logout.title'),
            content: t('notifications.auth.error_logout.content', [err.message]),
            type: 'error'
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

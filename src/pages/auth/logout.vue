<script lang="ts" setup>
import { LogoutRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { useAuthStore } from '~/store/auth';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import Footer from '~/components/partials/Footer.vue';
import { useNotificationsStore } from '~/store/notifications';

useHead({
    title: 'pages.auth.logout.title',
});
definePageMeta({
    title: 'pages.auth.logout.title',
    requiresAuth: true,
    authOnlyToken: true,
});

const { $grpc } = useNuxtApp();
const store = useAuthStore();
const notifications = useNotificationsStore();

const { t } = useI18n();

const accessToken = computed(() => store.$state.accessToken);

function redirect(): void {
    setTimeout(async () => {
        await navigateTo({ name: 'index' });
    }, 1500);
}

onBeforeMount(async () => {
    store.clear();

    if (!accessToken.value) {
        redirect();
        return;
    }

    try {
        await $grpc.getAuthClient()
            .logout(new LogoutRequest(), null);
    } catch (e) {
        const err = e as RpcError;
        store.loginStop(err.message);
        notifications.dispatchNotification({
            title: t('notifications.error_logout.title'),
            content: t('notifications.error_logout.content', [err.message]),
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

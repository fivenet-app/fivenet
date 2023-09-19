<script lang="ts" setup>
import Login from '~/components/auth/Login.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import Footer from '~/components/partials/Footer.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

useHead({
    title: 'pages.auth.login.title',
});
definePageMeta({
    title: 'pages.auth.login.title',
    requiresAuth: false,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { setAccessToken } = authStore;
const notifications = useNotificatorStore();
const route = useRoute();

const query = route.query;
// `t` and `exp` set, means social login was successful
if (query.t && query.t !== '' && query.exp) {
    setAccessToken(query.t as string, BigInt(query.exp as string));

    notifications.dispatchNotification({
        title: { key: 'notifications.auth.oauth2_login.success.title', parameters: [] },
        content: { key: 'notifications.auth.oauth2_login.success.content', parameters: [] },
        type: 'info',
    });

    await navigateTo({ name: 'auth-character-selector' });
    // `oauth2Login` can be `failed` (with `reason`)
} else if (query.oauth2Login && query.oauth2Login === 'failed') {
    const reason = query.reason ?? 'N/A';

    notifications.dispatchNotification({
        title: { key: 'notifications.auth.oauth2_login.failed.title', parameters: [] },
        content: { key: 'notifications.auth.oauth2_login.failed.content', parameters: [reason.toString()] },
        type: 'error',
    });
}
</script>

<template>
    <div class="h-full justify-between flex flex-col">
        <HeroFull>
            <ContentCenterWrapper>
                <Login />
            </ContentCenterWrapper>
        </HeroFull>
        <Footer />
    </div>
</template>

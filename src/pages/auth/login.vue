<script lang="ts" setup>
import Login from '~/components/auth/Login.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import Footer from '~/components/partials/Footer.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';

useHead({
    title: 'pages.auth.login.title',
});
definePageMeta({
    title: 'pages.auth.login.title',
    requiresAuth: false,
});

const authStore = useAuthStore();
const notifications = useNotificationsStore();
const route = useRoute();

const { setAccessToken } = authStore;

const { t } = useI18n();

const query = route.query;
if (query.t && query.t !== '' && query.exp) {
    setAccessToken(query.t as string, parseInt(query.exp as string));

    notifications.dispatchNotification({
        title: t('notifications.auth.oauth2_login.success.title'),
        content: t('notifications.auth.oauth2_login.success.content'),
        type: 'info',
    });

    await navigateTo({ name: 'auth-character-selector' });
} else if (query.oauth2Login && query.oauth2Login === 'failed') {
    const reason = query.reason ?? 'N/A';

    notifications.dispatchNotification({
        title: t('notifications.auth.oauth2_login.failed.title'),
        content: t('notifications.auth.oauth2_login.failed.content', [reason]),
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

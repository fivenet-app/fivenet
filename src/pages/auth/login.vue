<script lang="ts" setup>
import Login from '~/components/auth/Login.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
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

const token = route.query.token;
if (token && token !== "") {
    authStore.updateAccessToken(token as string);

    notifications.dispatchNotification({
        title: 'Successfully logged in',
        content: 'Successfully logged in using social login provider.',
        type: 'info',
    });
    await navigateTo({ name: 'auth-character-selector' });
}
</script>

<template>
    <HeroFull>
        <ContentCenterWrapper>
            <Login />
        </ContentCenterWrapper>
    </HeroFull>
</template>

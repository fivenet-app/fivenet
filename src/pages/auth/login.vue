<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import ForgotPasswordForm from '~/components/auth/ForgotPasswordForm.vue';
import LoginForm from '~/components/auth/LoginForm.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

useHead({
    title: 'components.auth.LoginForm.title',
});
definePageMeta({
    title: 'components.auth.LoginForm.title',
    layout: 'auth',
    requiresAuth: false,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { setAccessTokenExpiration } = authStore;
const { username } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const route = useRoute();

const showLogin = ref(true);

onMounted(async () => {
    const query = route.query;
    // `t` and `exp` set, means social login was successful
    if (query.u && query.u !== '' && query.exp && query.exp !== '') {
        console.info('Login: Got access token via query param (oauth2 login)');
        username.value = query.u as string;
        setAccessTokenExpiration(query.exp as string);

        notifications.add({
            title: { key: 'notifications.auth.oauth2_login.success.title', parameters: {} },
            description: { key: 'notifications.auth.oauth2_login.success.content', parameters: {} },
            type: NotificationType.INFO,
        });

        await navigateTo({ name: 'auth-character-selector' });
    } else if (query.oauth2Login && query.oauth2Login === 'failed') {
        // `oauth2Login` can be `failed` (with `reason`)
        const reason = query.reason ?? 'N/A';

        notifications.add({
            title: { key: 'notifications.auth.oauth2_login.failed.title', parameters: {} },
            description: { key: 'notifications.auth.oauth2_login.failed.content', parameters: { msg: reason.toString() } },
            type: NotificationType.ERROR,
        });
    }
});
</script>

<template>
    <UCard class="w-full max-w-sm bg-white/75 backdrop-blur dark:bg-white/5">
        <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

        <LoginForm v-if="showLogin" @toggle="showLogin = !showLogin" />
        <ForgotPasswordForm v-else @toggle="showLogin = !showLogin" />
    </UCard>
</template>

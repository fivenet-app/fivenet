<script lang="ts" setup>
import { type NavigationFailure } from 'vue-router';
import type { TypedRouteFromName } from '@typed-router';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import ForgotPasswordForm from '~/components/auth/ForgotPasswordForm.vue';
import LoginForm from '~/components/auth/LoginForm.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

useHead({
    title: 'components.auth.login.title',
});
definePageMeta({
    title: 'components.auth.login.title',
    layout: 'auth',
    requiresAuth: false,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { setAccessToken } = authStore;
const { accessToken } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const route = useRoute();

const showLogin = ref(true);

onMounted(async () => {
    const query = route.query;
    // `t` and `exp` set, means social login was successful
    if (query.t && query.t !== '' && query.exp && query.exp !== '') {
        console.info('Login: Got access token via query param (oauth2 login)');
        setAccessToken(query.t as string, BigInt(query.exp as string));

        notifications.dispatchNotification({
            title: { key: 'notifications.auth.oauth2_login.success.title', parameters: {} },
            content: { key: 'notifications.auth.oauth2_login.success.content', parameters: {} },
            type: 'info',
        });

        await navigateTo({ name: 'auth-character-selector' });
    } else if (query.oauth2Login && query.oauth2Login === 'failed') {
        // `oauth2Login` can be `failed` (with `reason`)
        const reason = query.reason ?? 'N/A';

        notifications.dispatchNotification({
            title: { key: 'notifications.auth.oauth2_login.failed.title', parameters: {} },
            content: { key: 'notifications.auth.oauth2_login.failed.content', parameters: { msg: reason.toString() } },
            type: 'error',
        });
    }
});

watch(accessToken, async (): Promise<NavigationFailure | TypedRouteFromName<'auth-character-selector'> | void | undefined> => {
    if (accessToken.value === null) {
        return;
    }

    return await navigateTo({
        name: 'auth-character-selector',
        query: route.query,
    });
});
</script>

<template>
    <UCard class="max-w-sm w-full bg-white/75 dark:bg-white/5 backdrop-blur">
        <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

        <LoginForm v-if="showLogin" @toggle="showLogin = !showLogin" />
        <ForgotPasswordForm v-else @toggle="showLogin = !showLogin" />
    </UCard>
</template>

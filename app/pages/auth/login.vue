<script lang="ts" setup>
import type { TabsItem } from '@nuxt/ui';
import ForgotPasswordForm from '~/components/auth/ForgotPasswordForm.vue';
import LoginForm from '~/components/auth/LoginForm.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { useAuthStore } from '~/stores/auth';
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

const { auth } = useAppConfig();

const { t } = useI18n();

const authStore = useAuthStore();
const { setUserToken } = authStore;
const { username } = storeToRefs(authStore);

const notifications = useNotificationsStore();

const logger = useLogger('🔑 Auth');

const items: TabsItem[] = [
    {
        slot: 'login' as const,
        label: t('components.auth.LoginForm.title'),
        icon: 'i-mdi-login',
        value: 'login',
    },
    {
        slot: 'forgotPassword' as const,
        label: t('components.auth.LoginForm.forgot_password'),
        icon: 'i-mdi-forgot-password',
        value: 'forgotPassword',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'login';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

onMounted(async () => {
    const query = route.query;
    // `u` set, means social login was successful
    if (query.u && query.u !== '') {
        logger.info('Got access token via query param (oauth2 login)');
        // Read and attempt to parse token from user token cookie
        const fivenetUsr = await cookieStore.get('fivenet_usr');
        if (!fivenetUsr?.value) return;

        setUserToken(fivenetUsr.value);
        username.value = query.u as string;

        // Remove the cookie after reading
        await cookieStore.delete('fivenet_usr');

        notifications.add({
            title: { key: 'notifications.auth.oauth2_login.success.title', parameters: {} },
            description: { key: 'notifications.auth.oauth2_login.success.content', parameters: {} },
            type: NotificationType.INFO,
        });

        await navigateTo('/auth/character-selector');
        return;
    }

    if (query.oauth2Login && query.oauth2Login === 'failed') {
        // `oauth2Login` can be `failed` (with an optional `reason` param explaining why)
        const reason = query.reason ?? 'N/A';

        notifications.add({
            title: { key: 'notifications.auth.oauth2_login.failed.title', parameters: {} },
            description: { key: 'notifications.auth.oauth2_login.failed.content', parameters: { msg: reason.toString() } },
            type: NotificationType.ERROR,
        });
    }
});

const canSubmit = ref(true);
</script>

<template>
    <UPageCard class="w-full max-w-md shrink-0 bg-white/75 backdrop-blur-sm dark:bg-white/5" :ui="{ header: 'w-full' }">
        <template #header>
            <div class="w-full space-y-2">
                <FiveNetLogo class="mx-auto h-auto w-20" />

                <h2 class="text-center text-3xl">
                    {{ $t('common.login') }}
                </h2>
            </div>
        </template>

        <div class="space-y-4">
            <UTabs v-model="selectedTab" class="w-full" :items="items">
                <template #login>
                    <LoginForm v-model="canSubmit" />
                </template>

                <template #forgotPassword>
                    <ForgotPasswordForm v-model="canSubmit" @toggle="selectedTab = 'login'" />
                </template>
            </UTabs>

            <div v-if="auth.signupEnabled" class="space-y-4">
                <USeparator orientation="horizontal" color="gray" />

                <UButton
                    block
                    color="neutral"
                    variant="outline"
                    trailing-icon="i-mdi-account-plus"
                    to="/auth/registration"
                    :disabled="!canSubmit"
                >
                    {{ $t('components.auth.LoginForm.register_account') }}
                </UButton>
            </div>
        </div>
    </UPageCard>
</template>

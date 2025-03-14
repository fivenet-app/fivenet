<script lang="ts" setup>
import ForgotPasswordForm from '~/components/auth/ForgotPasswordForm.vue';
import LoginForm from '~/components/auth/LoginForm.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { logger as authLogger, useAuthStore } from '~/stores/auth';
import { useNotificatorStore } from '~/stores/notificator';
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

const { login } = useAppConfig();

const { t } = useI18n();

const authStore = useAuthStore();
const { setAccessTokenExpiration } = authStore;
const { username } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const items = [
    {
        slot: 'login',
        label: t('components.auth.LoginForm.title'),
        icon: 'i-mdi-login',
    },
    {
        slot: 'forgotPassword',
        label: t('components.auth.LoginForm.forgot_password'),
        icon: 'i-mdi-forgot-password',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

onMounted(async () => {
    const query = route.query;
    // `t` and `exp` set, means social login was successful
    if (query.u && query.u !== '' && query.exp && query.exp !== '') {
        authLogger.info('Got access token via query param (oauth2 login)');
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

const canSubmit = ref(true);
</script>

<template>
    <UCard class="w-full max-w-md bg-white/75 backdrop-blur dark:bg-white/5">
        <div class="space-y-4">
            <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

            <UTabs v-model="selectedTab" :items="items" class="w-full">
                <template #login>
                    <LoginForm v-model="canSubmit" />
                </template>
                <template #forgotPassword>
                    <ForgotPasswordForm v-model="canSubmit" @toggle="selectedTab = 0" />
                </template>
            </UTabs>

            <div v-if="login.signupEnabled" class="space-y-4">
                <UDivider orientation="horizontal" />

                <UButton
                    block
                    color="gray"
                    trailing-icon="i-mdi-account-plus"
                    :to="{ name: 'auth-registration' }"
                    :disabled="!canSubmit"
                >
                    {{ $t('components.auth.LoginForm.register_account') }}
                </UButton>
            </div>
        </div>
    </UCard>
</template>

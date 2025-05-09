<script lang="ts" setup>
import type { ButtonColor, ButtonSize } from '#ui/types';
import '~/assets/css/herofull-pattern.css';
import { useAuthStore } from '~/stores/auth';

useHead({
    title: 'common.home',
});

definePageMeta({
    title: 'common.home',
    layout: 'landing',
    requiresAuth: false,
    showCookieOptions: true,
});

const { t } = useI18n();

const appConfig = useAppConfig();

const authStore = useAuthStore();
const { username } = storeToRefs(authStore);

const appVersion = APP_VERSION.split('-')[0];

const links = computed(() =>
    [
        username.value
            ? { label: t('common.overview'), icon: 'i-mdi-home', size: 'lg' as ButtonSize, to: '/overview' }
            : {
                  label: t('components.auth.LoginForm.title'),
                  icon: 'i-mdi-login',
                  size: 'lg' as ButtonSize,
                  to: '/auth/login',
              },
        appConfig.login.signupEnabled
            ? {
                  label: t('components.auth.RegistrationForm.title'),
                  trailingIcon: 'i-mdi-account-plus',
                  color: 'gray' as ButtonColor,
                  size: 'lg' as ButtonSize,
                  to: '/auth/registration',
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <div class="flex min-h-[calc(100dvh-(2*var(--header-height)))] flex-col">
        <div class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]" />

        <ULandingHero :title="$t('pages.index.welcome')" :description="$t('pages.index.subtext')" :links="links" class="flex-1">
            <template #headline>
                <UButton
                    color="gray"
                    :to="`https://github.com/fivenet-app/fivenet/v2025/releases/tag/${appVersion}`"
                    :external="true"
                    :label="$t('pages.index.whats_new_in', { version: appVersion })"
                    trailing-icon="i-mdi-arrow-right"
                    size="xs"
                    class="rounded-full"
                />
            </template>
        </ULandingHero>
    </div>
</template>

<style scoped>
.gradient {
    mask-image: radial-gradient(100% 100% at top, white, transparent);
}

.dark {
    .gradient {
        mask-image: radial-gradient(100% 100% at top, black, transparent);
    }
}
</style>

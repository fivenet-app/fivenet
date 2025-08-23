<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';
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

const { login } = useAppConfig();

const authStore = useAuthStore();
const { username } = storeToRefs(authStore);

const appVersion = APP_VERSION.split('-')[0];

const links = computed<ButtonProps[]>(() =>
    (
        [
            username.value
                ? { label: t('common.overview'), icon: 'i-mdi-home', size: 'lg', to: '/overview' }
                : {
                      label: t('components.auth.LoginForm.title'),
                      icon: 'i-mdi-login',
                      size: 'lg',
                      to: '/auth/login',
                  },
            login.signupEnabled
                ? {
                      label: t('components.auth.RegistrationForm.title'),
                      trailingIcon: 'i-mdi-account-plus',
                      color: 'neutral',
                      size: 'lg',
                      to: '/auth/registration',
                  }
                : undefined,
        ] as ButtonProps[]
    ).flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <div class="flex flex-1 flex-col">
        <div class="hero absolute inset-0 z-[-1] mask-[radial-gradient(100%_100%_at_top,white,transparent)]" />

        <div class="flex min-h-[calc(100dvh-var(--ui-header-height))] flex-col items-center justify-center">
            <UCard class="w-full max-w-4xl bg-white/75 backdrop-blur-sm dark:bg-white/5">
                <div class="space-y-4">
                    <UPageHero
                        :title="$t('pages.index.welcome')"
                        :description="$t('pages.index.subtext')"
                        :links="links"
                        :ui="{ wrapper: 'py-0 sm:py-0 md:py-0 relative', title: 'text-4xl' }"
                    >
                        <template #headline>
                            <UButton
                                class="rounded-full"
                                color="neutral"
                                :to="`https://github.com/fivenet-app/fivenet/releases/tag/${appVersion}`"
                                external
                                :label="$t('pages.index.whats_new_in', { version: appVersion })"
                                trailing-icon="i-mdi-arrow-right"
                                size="xs"
                            />
                        </template>
                    </UPageHero>
                </div>
            </UCard>
        </div>
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

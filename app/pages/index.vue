<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';
import type { ButtonProps } from '@nuxt/ui';
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

const { auth } = useAppConfig();

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
            auth.signupEnabled
                ? {
                      label: t('components.auth.registration_form.title'),
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
        <div class="flex min-h-[calc(100dvh-var(--ui-header-height))] flex-col items-center justify-center">
            <UCard
                class="mt-[calc(var(--page-content-bottom-offset)+4*var(--spacing))] mb-6 w-full max-w-4xl bg-white/75 backdrop-blur-sm dark:bg-white/5"
            >
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
                                variant="outline"
                                :to="`https://fivenet.app/changelog#${appVersion}`"
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

<script lang="ts" setup>
import type { Button } from '#ui/types';
import { useAuthStore } from '~/store/auth';
import '~/assets/css/herofull-pattern.css';

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

const links = computed(
    () =>
        [
            username.value
                ? { label: t('common.overview'), icon: 'i-mdi-home', size: 'lg', to: '/overview' }
                : {
                      label: t('components.auth.LoginForm.title'),
                      icon: 'i-mdi-login',
                      size: 'lg',
                      to: '/auth/login',
                  },
            appConfig.login.signupEnabled
                ? {
                      label: t('components.auth.RegistrationForm.title'),
                      trailingIcon: 'i-mdi-account-plus',
                      color: 'gray',
                      size: 'lg',
                      to: '/auth/registration',
                  }
                : undefined,
        ].flatMap((item) => (item !== undefined ? [item] : [])) as (Button & { click?: Function })[],
);
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero :title="$t('pages.index.welcome')" :description="$t('pages.index.subtext')" :links="links">
                <template #headline>
                    <UButton
                        color="gray"
                        :to="`https://github.com/fivenet-app/fivenet/releases/tag/${appVersion}`"
                        :external="true"
                        :label="$t('pages.index.whats_new_in', { version: appVersion })"
                        trailing-icon="i-mdi-arrow-right"
                        size="xs"
                        class="rounded-full"
                    />
                </template>
            </ULandingHero>
        </div>
    </div>
</template>

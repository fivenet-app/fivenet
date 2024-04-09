<script lang="ts" setup>
import type { Button } from '#ui/types';
import { useAuthStore } from '~/store/auth';
import { useSettingsStore } from '~/store/settings';

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
const { accessToken } = storeToRefs(authStore);

const router = useRouter();

onBeforeMount(async () => {
    if (accessToken.value) {
        // @ts-ignore the route should be valid, as we test it against a valid URL list
        const target = router.resolve(useSettingsStore().startpage ?? '/overview');
        return navigateTo(target);
    }
});

const appVersion = __APP_VERSION__.split('-')[0];

const links = computed(
    () =>
        [
            accessToken.value
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
                        :to="`https://github.com/galexrt/fivenet/releases/tag/${appVersion}`"
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

<script lang="ts" setup>
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
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero
                :title="$t('pages.index.welcome')"
                :description="$t('pages.index.subtext')"
                :links="
                    !accessToken
                        ? [
                              {
                                  label: $t('components.auth.LoginForm.title'),
                                  icon: 'i-mdi-login',
                                  size: 'lg',
                                  to: '/auth/login',
                              },
                              {
                                  label: $t('components.auth.RegistrationForm.title'),
                                  trailingIcon: 'i-mdi-account-plus',
                                  color: 'gray',
                                  size: 'lg',
                                  to: '/auth/registration',
                              },
                          ]
                        : [{ label: $t('common.overview'), icon: 'i-mdi-home', size: 'lg', to: '/overview' }]
                "
            >
                <template #headline>
                    <UButton
                        color="gray"
                        to="/about"
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

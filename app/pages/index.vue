<script lang="ts" setup>
import type { ButtonColor, ButtonSize } from '#ui/types';
import '~/assets/css/herofull-pattern.css';
import LanguageSwitcherModal from '~/components/partials/LanguageSwitcherModal.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
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

const modal = useModal();

const { login } = useAppConfig();

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
        login.signupEnabled
            ? {
                  label: t('components.auth.RegistrationForm.title'),
                  trailingIcon: 'i-mdi-account-plus',
                  color: 'white' as ButtonColor,
                  size: 'lg' as ButtonSize,
                  to: '/auth/registration',
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <div class="flex flex-1 flex-col">
        <div class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]" />

        <div
            class="flex max-h-[calc(100dvh-(var(--header-height)))] min-h-[calc(100dvh-(var(--header-height)))] flex-col items-center justify-center"
        >
            <div class="absolute top-4 z-10 flex gap-2">
                <UButton
                    color="white"
                    :label="$t('common.language')"
                    icon="i-mdi-translate"
                    @click="modal.open(LanguageSwitcherModal, {})"
                />
            </div>

            <UCard class="w-full max-w-4xl bg-white/75 backdrop-blur dark:bg-white/5">
                <div class="space-y-4">
                    <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

                    <ULandingHero
                        :title="$t('pages.index.welcome')"
                        :description="$t('pages.index.subtext')"
                        :links="links"
                        :ui="{ wrapper: 'py-0 sm:py-0 md:py-0 relative', title: 'text-4xl' }"
                    >
                        <template #headline>
                            <UButton
                                class="rounded-full"
                                color="gray"
                                :to="`https://github.com/fivenet-app/fivenet/releases/tag/${appVersion}`"
                                :external="true"
                                :label="$t('pages.index.whats_new_in', { version: appVersion })"
                                trailing-icon="i-mdi-arrow-right"
                                size="xs"
                            />
                        </template>
                    </ULandingHero>
                </div>
            </UCard>

            <UButton
                class="absolute bottom-4 z-10"
                icon="i-mdi-information-outline"
                :label="$t('pages.about.title')"
                to="/about"
                color="white"
            />
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

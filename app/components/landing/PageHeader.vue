<script setup lang="ts">
import LanguageSwitcherModal from '~/components/partials/LanguageSwitcherModal.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { useAuthStore } from '~/stores/auth';

const { t } = useI18n();

const appConfig = useAppConfig();

const authStore = useAuthStore();
const { username } = storeToRefs(authStore);

const { website } = useAppConfig();

const items = computed(() =>
    [
        !username.value
            ? {
                  label: t('common.home'),
                  icon: 'i-mdi-home',
                  to: '/',
              }
            : {
                  label: t('common.overview'),
                  icon: 'i-mdi-overview',
                  to: '/overview',
              },
        website.statsPage
            ? {
                  label: t('pages.stats.title'),
                  icon: 'i-mdi-analytics',
                  to: '/stats',
              }
            : undefined,
        {
            label: t('common.about'),
            icon: 'i-mdi-information',
            to: '/about',
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const overlay = useOverlay();
const langSwitcherModal = overlay.create(LanguageSwitcherModal, {});
</script>

<template>
    <UHeader>
        <template #title>
            <FiveNetLogo class="h-10 w-auto" />
        </template>

        <UNavigationMenu :items="items" />

        <template #right>
            <UButton :label="$t('common.language')" icon="i-mdi-translate" color="neutral" @click="langSwitcherModal.open()" />

            <template v-if="!username">
                <UButton :label="$t('components.auth.LoginForm.title')" icon="i-mdi-login" color="neutral" to="/auth/login" />
                <UButton
                    v-if="appConfig.login.signupEnabled"
                    :label="$t('components.auth.RegistrationForm.title')"
                    icon="i-mdi-account-plus"
                    trailing
                    color="neutral"
                    to="/auth/registration"
                    class="hidden lg:flex"
                />
            </template>
            <template v-else>
                <UButton :label="$t('common.overview')" icon="i-mdi-home" to="/overview" />
            </template>
        </template>
    </UHeader>
</template>

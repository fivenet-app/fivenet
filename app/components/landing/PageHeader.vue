<script setup lang="ts">
import LanguageSwitcherModal from '~/components/partials/LanguageSwitcherModal.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { useAuthStore } from '~/stores/auth';

const { t } = useI18n();

const { login } = useAppConfig();

const authStore = useAuthStore();
const { username } = storeToRefs(authStore);

const { website } = useAppConfig();

const links = computed(() =>
    [
        !username.value
            ? {
                  label: t('common.home'),
                  to: '/',
              }
            : {
                  label: t('common.overview'),
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
            to: '/about',
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const modal = useModal();
</script>

<template>
    <UHeader :links="links" :ui="{ logo: 'inline-flex items-center gap-2' }">
        <template #logo>
            <FiveNetLogo class="h-10 w-auto" />

            <span class="text-xl font-semibold text-gray-900 dark:text-white">FiveNet</span>
        </template>

        <template #right>
            <UButton
                :label="$t('common.language')"
                icon="i-mdi-translate"
                color="gray"
                @click="modal.open(LanguageSwitcherModal, {})"
            />

            <template v-if="!username">
                <UButton :label="$t('components.auth.LoginForm.title')" icon="i-mdi-login" color="gray" to="/auth/login" />

                <UButton
                    v-if="login.signupEnabled"
                    class="hidden lg:flex"
                    :label="$t('components.auth.RegistrationForm.title')"
                    icon="i-mdi-account-plus"
                    trailing
                    color="black"
                    to="/auth/registration"
                />
            </template>
        </template>
    </UHeader>
</template>

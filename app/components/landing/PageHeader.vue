<script setup lang="ts">
import { de, en } from '@nuxt/ui/locale';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { useAuthStore } from '~/stores/auth';

const { t } = useI18n();

const { login } = useAppConfig();

const authStore = useAuthStore();
const { username } = storeToRefs(authStore);

const { website } = useAppConfig();

const items = computed(() =>
    [
        !username.value
            ? {
                  label: t('common.home'),
                  to: '/',
                  icon: 'i-mdi-home',
              }
            : {
                  label: t('common.overview'),
                  to: '/overview',
                  icon: 'i-mdi-view-dashboard',
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
            icon: 'i-mdi-information',
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const settingsStore = useSettingsStore();
const { locale: userLocale } = storeToRefs(settingsStore);

const { locale, setLocale } = useI18n();
</script>

<template>
    <UHeader :ui="{ title: 'inline-flex items-center gap-2' }">
        <template #title>
            <FiveNetLogo class="h-10 w-auto" />

            <span class="text-xl font-semibold text-highlighted">FiveNet</span>
        </template>

        <UNavigationMenu :items="items" />

        <template #right>
            <ULocaleSelect
                v-model="locale"
                :locales="[en, de]"
                @update:model-value="
                    ($event) => {
                        let l = $event as typeof userLocale;
                        if (!l) l = 'en';
                        setLocale(l);
                        userLocale = l as typeof userLocale;
                    }
                "
            />

            <template v-if="!username">
                <UButton :label="$t('components.auth.LoginForm.title')" icon="i-mdi-login" to="/auth/login" />

                <UButton
                    v-if="login.signupEnabled"
                    class="hidden lg:flex"
                    :label="$t('components.auth.RegistrationForm.title')"
                    icon="i-mdi-account-plus"
                    trailing
                    color="neutral"
                    to="/auth/registration"
                />
            </template>
        </template>
    </UHeader>
</template>

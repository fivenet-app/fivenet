<script setup lang="ts">
import { de, en } from '@nuxt/ui/locale';
import BannerMessage from '~/components/partials/BannerMessage.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

const { t } = useI18n();

const { auth, website, system } = useAppConfig();

const { username } = useAuth();

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

const { currentLocale, setUserLocale } = useAppLocale();

async function changeLocale(newLocale: string) {
    await setUserLocale(newLocale);
}

const currentRoute = useRoute();
const isAuthPage = computed(() => currentRoute.path.startsWith('/auth'));
</script>

<template>
    <UHeader :ui="{ title: 'inline-flex items-center gap-2' }">
        <template #title>
            <FiveNetLogo class="h-10 w-auto" />

            <span class="text-xl font-semibold text-highlighted">FiveNet</span>
        </template>

        <UNavigationMenu :items="items" />

        <template #right>
            <ULocaleSelect v-model="currentLocale" :locales="[en, de]" @update:model-value="($event) => changeLocale($event)" />

            <template v-if="!username && !isAuthPage">
                <UButton :label="$t('components.auth.LoginForm.title')" icon="i-mdi-login" to="/auth/login" />

                <UButton
                    v-if="auth.signupEnabled"
                    class="hidden lg:flex"
                    :label="$t('components.auth.registration_form.title')"
                    icon="i-mdi-account-plus"
                    trailing
                    color="neutral"
                    to="/auth/registration"
                />
            </template>
        </template>

        <template #body>
            <UNavigationMenu :items="items" orientation="vertical" class="-mx-2.5" />
        </template>

        <template #bottom>
            <BannerMessage v-if="system.bannerMessageEnabled && system.bannerMessage" :message="system.bannerMessage" />
        </template>
    </UHeader>
</template>

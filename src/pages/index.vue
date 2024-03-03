<script lang="ts" setup>
import HeroPage from '~/components/partials/HeroPage.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { useAuthStore } from '~/store/auth';
import { useSettingsStore } from '~/store/settings';

useHead({
    title: 'common.home',
});
definePageMeta({
    title: 'common.home',
    requiresAuth: false,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);

onBeforeMount(async () => {
    if (accessToken.value) {
        // @ts-ignore the route should be valid, as we test it against a valid URL list
        const target = useRouter().resolve(useSettingsStore().startpage ?? '/overview');
        return navigateTo(target);
    }
});
</script>

<template>
    <HeroPage>
        <template #default>
            <FiveNetLogo class="mx-auto mb-2 h-36 w-auto" />

            <h1 class="text-5xl font-bold tracking-tight text-neutral sm:text-6xl">
                {{ $t('pages.index.welcome') }}
            </h1>

            <p class="mt-6 text-lg leading-8 text-neutral">
                {{ $t('pages.index.subtext') }}
            </p>

            <div class="mt-4 flex items-center justify-center gap-x-6">
                <NuxtLink
                    :to="{ name: 'auth-login' }"
                    class="w-48 max-w-96 rounded-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                >
                    {{ $t('components.auth.login.title') }}
                </NuxtLink>
                <NuxtLink
                    :to="{ name: 'auth-registration' }"
                    class="w-48 max-w-96 rounded-md bg-secondary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-secondary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-secondary-500"
                >
                    {{ $t('components.auth.registration_form.title') }}
                </NuxtLink>
            </div>
        </template>
    </HeroPage>
</template>

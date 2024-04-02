<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
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
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero
                :title="$t('pages.index.welcome')"
                :description="$t('pages.index.subtext')"
                :links="[
                    { label: $t('components.auth.login.title'), icon: 'i-mdi-login', size: 'lg', to: '/auth/login' },
                    {
                        label: $t('components.auth.registration_form.title'),
                        trailingIcon: 'i-mdi-account-plus',
                        color: 'gray',
                        size: 'lg',
                        to: '/auth/registration',
                    },
                ]"
            >
                <template #headline>
                    <FiveNetLogo class="mx-auto h-36 w-auto" />
                </template>
            </ULandingHero>
        </div>
    </div>
</template>

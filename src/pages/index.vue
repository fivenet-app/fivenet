<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import PageFooter from '~/components/partials/PageFooter.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
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
        // @ts-expect-error the route should be valid, as we test it against a valid URL list
        const target = useRouter().resolve(useSettingsStore().startpage ?? '/overview');
        return navigateTo(target);
    }
});
</script>

<template>
    <div class="flex h-full flex-col justify-between">
        <HeroFull>
            <ContentCenterWrapper class="mx-auto max-w-2xl text-center">
                <div class="px-5 sm:px-0">
                    <FiveNetLogo class="mx-auto mb-2 h-auto w-36" />

                    <h1 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                        {{ $t('pages.index.welcome') }}
                    </h1>

                    <p v-t="'pages.index.subtext'" class="mt-6 text-lg leading-8 text-neutral"></p>
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
                </div>
            </ContentCenterWrapper>
        </HeroFull>
        <PageFooter />
    </div>
</template>

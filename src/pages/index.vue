<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import Footer from '~/components/partials/Footer.vue';
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
        const target = useRouter().resolve(useSettingsStore().startpage);
        return navigateTo(target);
    }
});
</script>

<template>
    <div class="h-full justify-between flex flex-col">
        <HeroFull>
            <ContentCenterWrapper class="max-w-2xl mx-auto text-center">
                <div class="sm:px-0 px-5">
                    <FiveNetLogo class="h-auto mx-auto mb-2 w-36" />

                    <h1 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                        {{ $t('pages.index.welcome') }}
                    </h1>

                    <p v-t="'pages.index.subtext'" class="mt-6 text-lg leading-8 text-neutral"></p>
                    <div class="flex items-center justify-center mt-4 gap-x-6">
                        <NuxtLink
                            :to="{ name: 'auth-login' }"
                            class="w-48 max-w-96 rounded-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        >
                            {{ $t('pages.auth.login.menu_item') }}
                        </NuxtLink>
                    </div>
                </div>
            </ContentCenterWrapper>
        </HeroFull>
        <Footer />
    </div>
</template>

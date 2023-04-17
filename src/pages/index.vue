<script lang="ts" setup>
import { computed } from 'vue';
import { useAuthStore } from '~/store/auth';
import Footer from '~/components/partials/Footer.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';

useHead({
    title: 'Home',
});
definePageMeta({
    title: 'Home',
    requiresAuth: false,
});

const store = useAuthStore();

const accessToken = computed(() => store.$state.accessToken);
</script>

<template>
    <div class="h-full justify-between flex flex-col">
        <HeroFull>
            <ContentCenterWrapper class="max-w-2xl mx-auto text-center">
                <div class="sm:px-0 px-5">
                    <img class="h-auto mx-auto mb-2 w-36" src="/images/logo.png" alt="FiveNet Logo" />
                    <h1 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                        {{ $t('pages.index.welcome') }}
                    </h1>
                    <p v-t="'pages.index.subtext'" class="mt-6 text-lg leading-8 text-neutral"></p>
                    <div class="flex items-center justify-center mt-4 gap-x-6">
                        <NuxtLink v-if="accessToken" :to="{ name: 'overview' }"
                            class="rounded-md w-32 bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                            {{ $t('overview') }}
                        </NuxtLink>
                        <NuxtLink v-else :to="{ name: 'auth-login' }"
                            class="rounded-md w-20 bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                            {{ $t('login') }}
                        </NuxtLink>
                    </div>
                </div>
            </ContentCenterWrapper>
        </HeroFull>
        <Footer />
    </div>
</template>

<script lang="ts" setup>
import HeroPage from '~/components/partials/HeroPage.vue';
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'pages.notfound.title',
});
definePageMeta({
    title: 'pages.notfound.title',
    requiresAuth: false,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);
</script>

<template>
    <HeroPage>
        <template #default>
            <div class="px-5 sm:px-0">
                <h1 class="text-5xl font-bold tracking-tight text-neutral sm:text-6xl">
                    {{ $t('pages.notfound.page_not_found') }}
                </h1>

                <h2 class="text-4xl text-neutral">
                    {{ $t('pages.notfound.error') }}
                </h2>

                <p class="mt-6 text-lg leading-8 text-neutral">
                    {{ $t('pages.notfound.fun_error') }}
                </p>

                <div class="inline-flex items-center justify-center gap-x-6">
                    <NuxtLink
                        v-if="accessToken"
                        :to="{ name: 'overview' }"
                        class="w-48 max-w-96 rounded-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    >
                        {{ $t('common.overview') }}
                    </NuxtLink>
                    <NuxtLink
                        v-else
                        :to="{ name: 'index' }"
                        class="w-48 max-w-96 rounded-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    >
                        {{ $t('common.home') }}
                    </NuxtLink>
                </div>
            </div>
        </template>
    </HeroPage>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'pages.notfound.title',
});
definePageMeta({
    title: 'pages.notfound.title',
    layout: 'landing',
    requiresAuth: false,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero
                :title="$t('pages.notfound.page_not_found')"
                :description="$t('pages.notfound.fun_error')"
                :links="
                    accessToken
                        ? [
                              {
                                  label: $t('common.overview'),
                                  trailingIcon: 'i-mdi-account-plus',
                                  color: 'gray',
                                  size: 'lg',
                                  to: '/overview',
                              },
                          ]
                        : [{ label: $t('common.home'), icon: 'i-mdi-login', size: 'lg', to: '/' }]
                "
            >
                <template #headline>
                    <UBadge color="gray" variant="solid" size="lg">{{ $t('pages.notfound.error') }}</UBadge>
                </template>
            </ULandingHero>
        </div>
    </div>
</template>

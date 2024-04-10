<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';

useHead({
    title: 'common.home',
});
definePageMeta({
    title: 'common.home',
    layout: 'landing',
    requiresAuth: false,
    showCookieOptions: true,
});

const { t } = useI18n();

const links = [{ label: t('common.docs'), icon: 'i-mdi-book-open-variant-outline', size: 'lg', to: '/docs' }];

const appVersion = __APP_VERSION__.split('-')[0];
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero :title="$t('pages.index.welcome')" :description="$t('pages.index.subtext')" :links="links">
                <template #headline>
                    <UButton
                        color="gray"
                        :to="`https://github.com/galexrt/fivenet/releases/tag/${appVersion}`"
                        :external="true"
                        :label="$t('pages.index.whats_new_in', { version: appVersion })"
                        trailing-icon="i-mdi-arrow-right"
                        size="xs"
                        class="rounded-full"
                    />
                </template>
            </ULandingHero>
        </div>
    </div>
</template>

<script setup lang="ts">
import BannerMessage from '~/components/partials/BannerMessage.vue';
import PageFooter from '~/components/partials/PageFooter.vue';

const { system } = useAppConfig();
</script>

<!-- eslint-disable tailwindcss/no-custom-classname -->
<template>
    <div class="h-dvh">
        <div class="hero absolute inset-0 z-[-1]" />

        <div class="flex min-h-dvh flex-col items-center justify-center">
            <UButton class="absolute top-4 z-10" icon="i-mdi-home" :label="$t('common.home')" to="/" color="neutral" />

            <slot />
        </div>

        <BannerMessage v-if="system.bannerMessageEnabled && system.bannerMessage" :message="system.bannerMessage" />

        <ClientOnly>
            <LazyPartialsEventsLayer />
        </ClientOnly>

        <PageFooter />
    </div>
</template>

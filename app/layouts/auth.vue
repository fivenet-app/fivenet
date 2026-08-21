<script setup lang="ts">
import '~/assets/css/herofull-pattern.css';
import BannerMessage from '~/components/partials/BannerMessage.vue';

const { system } = useAppConfig();
</script>

<!-- eslint-disable tailwindcss/no-custom-classname -->
<template>
    <div class="flex h-screen items-center justify-center px-4">
        <div class="hero absolute inset-0 z-[-1]" />

        <UTooltip :text="$t('common.back')">
            <UButton
                icon="i-mdi-chevron-left"
                to="/"
                size="xl"
                color="neutral"
                variant="subtle"
                class="absolute top-8 left-8 z-10 mt-(--page-content-bottom-offset) rounded-full"
            />
        </UTooltip>

        <slot />

        <div
            v-if="system.bannerMessageEnabled && system.bannerMessage"
            class="absolute inset-x-0 top-0 z-[49] flex max-h-21 flex-col gap-0 print:hidden"
        >
            <BannerMessage :message="system.bannerMessage" />
        </div>
    </div>
</template>

<script setup lang="ts">
import '~/assets/css/herofull-pattern.css';
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
                class="absolute top-8 left-8 z-10 rounded-full"
            />
        </UTooltip>

        <slot />
    </div>
</template>

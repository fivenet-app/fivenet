<script setup lang="ts">
import PageFooter from '~/components/partials/PageFooter.vue';
</script>

<!-- eslint-disable tailwindcss/no-custom-classname -->
<template>
    <div class="h-dscreen">
        <div class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]" />

        <div class="flex min-h-[calc(100dvh-var(--header-height))] flex-col items-center justify-center">
            <UButton class="absolute top-4 z-10" icon="i-mdi-home" :label="$t('common.home')" to="/" color="black" />

            <slot />
        </div>

        <PageFooter />
    </div>
</template>

<script setup lang="ts">
import PageFooter from '~/components/partials/PageFooter.vue';

useHead({
    bodyAttrs: {
        class: 'dark:bg-gray-950',
    },
});
</script>

<!-- eslint-disable tailwindcss/no-custom-classname -->
<template>
    <div class="hero h-screen">
        <div class="flex h-full flex-col items-center justify-center">
            <div class="gradient" />

            <UButton icon="i-mdi-home" :label="$t('common.home')" to="/" color="black" class="absolute top-4" />

            <slot />
        </div>

        <PageFooter />
    </div>
</template>

<style scoped>
.gradient {
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: radial-gradient(50% 50% at 50% 50%, rgb(var(--color-primary-500) / 0.25) 0, #fff 100%);
}

.dark {
    .gradient {
        background: radial-gradient(50% 50% at 50% 50%, rgb(var(--color-primary-400) / 0.1) 0, rgb(var(--color-gray-950)) 100%);
    }
}

.dark {
    .overlay {
        background-image: linear-gradient(to right, rgb(var(--color-gray-900)) 0.5px, transparent 0.5px),
            linear-gradient(to bottom, rgb(var(--color-gray-900)) 0.5px, transparent 0.5px);
    }
}
</style>

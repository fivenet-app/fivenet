<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';
</script>

<template>
    <div class="isolate bg-base-900 hero h-full">
        <div class="hero-overlay h-full">
            <slot />
        </div>
    </div>
</template>

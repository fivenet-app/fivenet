<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';
</script>

<template>
    <div class="isolate h-full hero bg-base-900">
        <div class="h-full hero-overlay">
            <slot />
        </div>
    </div>
</template>

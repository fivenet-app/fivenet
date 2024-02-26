<script lang="ts" setup>
import '~/assets/css/herofull-pattern.css';
</script>

<template>
    <div class="h-full hero bg-base-900">
        <div class="hero-overlay h-full">
            <slot />
        </div>
    </div>
</template>

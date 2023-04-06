<script setup lang="ts">
import ContentWrapper from './ContentWrapper.vue';
</script>

<template>
    <div class="flex overflow-hidden justify-center items-center align-middle h-full">
        <main>
            <div class="mx-auto max-w-7xl">
                <slot></slot>
            </div>
        </main>
    </div>
</template>

<script setup lang="ts">
import ContentWrapper from './ContentWrapper.vue';
</script>

<template>
    <div class="flex overflow-hidden justify-center items-center align-middle h-full">
        <ContentWrapper class="flex-1">
            <slot></slot>
        </ContentWrapper>
    </div>
</template>

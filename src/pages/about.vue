<script lang="ts" setup>
import ContentCenterWrapper from '../components/partials/ContentCenterWrapper.vue';

useHead({
    title: 'About',
});
definePageMeta({
    title: 'About',
    requiresAuth: false,
});
</script>

<template>
    <ContentCenterWrapper>
        <h1 class="text-neutral text-center text-7xl">About</h1>
    </ContentCenterWrapper>
</template>

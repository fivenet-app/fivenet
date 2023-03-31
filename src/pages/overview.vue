<script lang="ts" setup>
import ContentCenterWrapper from '../components/partials/ContentCenterWrapper.vue';
import FeatureCards from '../components/FeatureCards.vue';

useHead({
    title: 'Overview',
});
definePageMeta({
    requiresAuth: true,
    permission: 'Overview.View',
});
</script>

<template>
    <ContentCenterWrapper>
        <FeatureCards />
    </ContentCenterWrapper>
</template>

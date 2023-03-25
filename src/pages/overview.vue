<script lang="ts" setup>
import ContentCenterWrapper from '../components/partials/ContentCenterWrapper.vue';
import FeatureCards from '../components/FeatureCards.vue';
</script>

<route lang="json">
{
    "name": "Overview",
    "meta": {
        "requiresAuth": true,
        "permission": "Overview.View"
    }
}
</route>

<template>
    <ContentCenterWrapper>
        <FeatureCards />
    </ContentCenterWrapper>
</template>

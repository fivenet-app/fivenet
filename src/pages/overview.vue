<script lang="ts" setup>
import ContentWrapper from '../components/partials/ContentWrapper.vue';
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
    <ContentWrapper>
        <FeatureCards />
    </ContentWrapper>
</template>

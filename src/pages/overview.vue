<script lang="ts" setup>
import Footer from '../components/partials/Footer.vue';
import FeatureCards from '../components/FeatureCards.vue';
import ContentWrapper from '../components/partials/ContentWrapper.vue';
import NavPageHeader from '../components/partials/NavPageHeader.vue';
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
    <NavPageHeader title="Overview" />
    <ContentWrapper>
        <FeatureCards />
    </ContentWrapper>
    <Footer />
</template>

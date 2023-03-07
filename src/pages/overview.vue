<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../components/partials/Navbar.vue';
import Footer from '../components/partials/Footer.vue';
import FeatureCards from '../components/FeatureCards.vue';
import ContentWrapper from '../components/partials/ContentWrapper.vue';
import NavPageHeader from '../components/partials/NavPageHeader.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        FeatureCards,
        ContentWrapper,
        NavPageHeader,
    },
});
</script>

<route lang="json">
{
    "name": "overview",
    "meta": {
        "requiresAuth": true,
        "permission": "overview-view",
        "breadCrumbs": []
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Overview" />
    <ContentWrapper>
        <FeatureCards />
    </ContentWrapper>
    <Footer />
</template>

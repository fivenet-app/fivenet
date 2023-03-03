<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import FeatureCards from '../components/FeatureCards.vue';
import ContentWrapper from '../components/ContentWrapper.vue';
import NavPageHeader from '../components/NavPageHeader.vue';

export default defineComponent({
    components: {
    Navbar,
    Footer,
    FeatureCards,
    ContentWrapper,
    NavPageHeader
},
});
</script>

<route lang="json">
{
    "name": "overview",
    "meta": {
        "requiresAuth": true,
        "permission": "overview-view"
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

<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/partials/Navbar.vue';
import Footer from '../../components/partials/Footer.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';

export default defineComponent({
    components: {
    Navbar,
    Footer,
    ContentWrapper,
    NavPageHeader
},
});
</script>

<route lang="json">
{
    "name": "dispatches",
    "meta": {
        "requiresAuth": true,
        "permission": "dispatches-view",
        "breadCrumbs": [
            { "name": "Dispatches", "href": "/dispatches" }
        ]
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Dispatches" />
    <ContentWrapper>
        <h1 class="text-lg font-semibold leading-6 text-white">DISPATCH CENTER AND MANAGEMENT SYSTEM</h1>
    </ContentWrapper>
    <Footer />
</template>

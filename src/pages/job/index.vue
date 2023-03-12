<script lang="ts">
import Navbar from '../../components/partials/Navbar.vue';
import Footer from '../../components/partials/Footer.vue';

import { defineComponent } from 'vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';

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
    "name": "Job",
    "meta": {
        "requiresAuth": true,
        "permission": "jobs-view",
        "breadCrumbs": [
            { "name": "Job", "href": "/job" }
        ]
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Job" />
    <ContentWrapper>
        <h1 class="text-lg font-semibold leading-6 text-white">EMPLOYEE INFO AND MANAGEMENT SYSTEM</h1>
    </ContentWrapper>
    <Footer />
</template>

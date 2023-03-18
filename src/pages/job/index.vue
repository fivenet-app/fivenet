<script lang="ts">
import Footer from '../../components/partials/Footer.vue';

import { defineComponent } from 'vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';

export default defineComponent({
    components: {
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
        "permission": "Jobs.View",
        "breadCrumbs": [
            { "name": "Job", "href": "/job" }
        ]
    }
}
</route>

<template>
    <NavPageHeader title="Job" />
    <ContentWrapper>
        <h1 class="text-lg font-semibold leading-6 text-white">EMPLOYEE INFO AND MANAGEMENT SYSTEM</h1>
    </ContentWrapper>
    <Footer />
</template>

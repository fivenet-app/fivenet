<script lang="ts">
import Navbar from '../../components/partials/Navbar.vue';
import Footer from '../../components/partials/Footer.vue';

import { defineComponent } from 'vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import Breadcrumbs from '../../components/partials/Breadcrumbs.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';

export default defineComponent({
    components: {
    Navbar,
    Footer,
    ContentWrapper,
    Breadcrumbs,
    NavPageHeader
},
});
</script>

<route lang="json">
{
    "name": "jobs",
    "meta": {
        "requiresAuth": true,
        "permission": "job-view"
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Job" />
    <ContentWrapper>
        <h1>JOBS - EMPLOYEE MANAGEMENT SYSTEM</h1>
    </ContentWrapper>
    <Footer />
</template>

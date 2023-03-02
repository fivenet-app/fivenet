<script lang="ts">
import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';

import { defineComponent } from 'vue';
import ContentWrapper from '../../components/ContentWrapper.vue';
import Breadcrumbs from '../../components/Breadcrumbs.vue';
import NavPageHeader from '../../components/NavPageHeader.vue';

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
        "requiresAuth": true
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

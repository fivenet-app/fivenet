<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';
import CitizensList from '../../components/CitizensList.vue';
import ContentWrapper from '../../components/ContentWrapper.vue';
import Breadcrumbs from '../../components/Breadcrumbs.vue';
import NavPageHeader from '../../components/NavPageHeader.vue';

export default defineComponent({
    components: {
    Navbar,
    Footer,
    CitizensList,
    ContentWrapper,
    Breadcrumbs,
    NavPageHeader
},
});
</script>

<route lang="json">
{
    "name": "users",
    "meta": {
        "requiresAuth": true
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Citizens" />
    <ContentWrapper>
        <CitizensList />
    </ContentWrapper>
    <Footer />
</template>

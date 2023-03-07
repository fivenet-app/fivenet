<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../../components/partials/Navbar.vue';
import Footer from '../../components/partials/Footer.vue';
import CitizensList from '../../components/citizens/CitizensList.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import Breadcrumbs from '../../components/partials/Breadcrumbs.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';

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
        "requiresAuth": true,
        "permission": "users-findusers",
        "breadCrumbs": [
            { "name": "Citizens", "href": "/citizens" }
        ]
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

<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../../components/partials/Navbar.vue';
import Footer from '../../components/partials/Footer.vue';
import CitizensList from '../../components/citizens/CitizensList.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';

export default defineComponent({
    components: {
    Navbar,
    Footer,
    CitizensList,
    ContentWrapper,
    NavPageHeader
},
});
</script>

<route lang="json">
{
    "name": "Citizens",
    "meta": {
        "requiresAuth": true,
        "permission": "citizenstoreservice-findusers",
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

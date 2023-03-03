<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';
import DocumentsList from '../../components/DocumentsList.vue';
import ContentWrapper from '../../components/ContentWrapper.vue';
import NavPageHeader from '../../components/NavPageHeader.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        DocumentsList,
        ContentWrapper,
        NavPageHeader,
    },
});
</script>

<route lang="json">
{
    "name": "documents",
    "meta": {
        "requiresAuth": true,
        "permission": "documents-view"
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <DocumentsList />
    </ContentWrapper>
    <Footer />
</template>

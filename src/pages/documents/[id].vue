<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/partials/Navbar.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import Footer from '../../components/partials/Footer.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        ContentWrapper,
        NavPageHeader,
    },
});
</script>

<route lang="json">
{
    "name": "Documents: Info",
    "meta": {
        "requiresAuth": true,
        "permission": "documents-createdocument",
        "breadCrumbs": [
            { "name": "Documents", "href": "/documents" },
            { "name": "Documents", "href": "/documents" }
        ]
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Document" />
    <ContentWrapper>
        <h1 class="text-lg text-white">Document</h1>
    </ContentWrapper>
    <Footer />
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/partials/Navbar.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import Footer from '../../components/partials/Footer.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import { Document } from '@arpanet/gen/documents/documents_pb';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        ContentWrapper,
        NavPageHeader,
    },
    data() {
        return {
            document: undefined as undefined | Document,
        };
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
            { "name": "Document Info: ...", "href": "/documents" }
        ]
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Document" />
    <ContentWrapper>
        <h1 class="text-lg text-white">Document</h1>
        {{ document }}
    </ContentWrapper>
    <Footer />
</template>

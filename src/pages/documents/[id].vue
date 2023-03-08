<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/partials/Navbar.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import Footer from '../../components/partials/Footer.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentView from '../../components/documents/DocumentView.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        ContentWrapper,
        NavPageHeader,
        DocumentView
    },
    data() {
        return {
            documentID: 0 as number,
        };
    },
    beforeMount() {
        this.documentID = +this.$route.params.id;
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
        <DocumentView v-if="documentID > 0" :documentID="documentID" />
    </ContentWrapper>
    <Footer />
</template>

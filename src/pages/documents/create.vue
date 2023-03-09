<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/partials/Navbar.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentEditor from '../../components/documents/DocumentEditor.vue';
import Footer from '../../components/partials/Footer.vue';

export default defineComponent({
    name: 'app',
    components: {
        Navbar,
        Footer,
        ContentWrapper,
        NavPageHeader,
        DocumentEditor,
    },
});
</script>

<route lang="json">
{
    "name": "Documents: Create",
    "meta": {
        "requiresAuth": true,
        "permission": "documents-createdocument",
        "breadCrumbs": [
            { "name": "Documents", "href": "/documents" },
            { "name": "Create Document", "href": "#" }
        ]
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <DocumentEditor />
    </ContentWrapper>
    <Footer />
</template>

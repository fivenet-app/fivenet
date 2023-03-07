<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/partials/Navbar.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import Editor from '../../components/documents/Editor.vue';
import Footer from '../../components/partials/Footer.vue';

export default defineComponent({
    name: 'app',
    components: {
        Editor
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
            { "name": "Documents", "href": "/documents" }
        ]
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <Editor />
    </ContentWrapper>
    <Footer />
</template>

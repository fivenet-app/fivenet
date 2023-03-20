<script setup lang="ts">
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentEditor from '../../components/documents/DocumentEditor.vue';
import Footer from '../../components/partials/Footer.vue';
</script>

<route lang="json">
{
    "name": "Documents: Create",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.CreateDocument",
        "breadCrumbs": [
            { "name": "Documents", "href": "/documents" },
            { "name": "Create Document", "href": "#" }
        ]
    }
}
</route>

<template>
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <DocumentEditor :id="9 /* TODO: Remove this. It's only for testing rn */" />
    </ContentWrapper>
    <Footer />
</template>

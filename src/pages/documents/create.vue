<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentEditor from '../../components/documents/DocumentEditor.vue';
import Footer from '../../components/partials/Footer.vue';
</script>

<route lang="json">
{
    "name": "Documents: Create",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.CreateDocument"
    }
}
</route>

<template>
    <ContentWrapper>
        <DocumentEditor />
    </ContentWrapper>
    <Footer />
</template>

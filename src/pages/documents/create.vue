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
        "permission": "DocStoreService.CreateDocument"
    }
}
</route>

<template>
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <DocumentEditor />
    </ContentWrapper>
    <Footer />
</template>

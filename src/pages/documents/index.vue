<script setup lang="ts">
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentsList from '../../components/documents/DocumentsList.vue';
import Footer from '../../components/partials/Footer.vue';
</script>

<route lang="json">
{
    "name": "Documents",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.FindDocuments"
    }
}
</route>

<template>
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <DocumentsList />
    </ContentWrapper>
    <Footer />
</template>

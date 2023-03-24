<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentsList from '../../components/documents/DocumentsList.vue';
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
    <ContentWrapper>
        <DocumentsList />
    </ContentWrapper>
</template>

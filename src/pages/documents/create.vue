<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentEditor from '../../components/documents/DocumentEditor.vue';
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
</template>

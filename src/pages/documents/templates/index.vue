<script lang="ts" setup>
import ContentWrapper from '../../../components/partials/ContentWrapper.vue';
import TemplatesList from '../../../components/documents/TemplatesList.vue';
</script>

<route lang="json">
{
    "name": "Documents: Templates",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.FindDocuments"
    }
}
</route>

<template>
    <ContentWrapper>
        <TemplatesList />
    </ContentWrapper>
</template>

<script lang="ts" setup>
import ContentWrapper from '../../../components/partials/ContentWrapper.vue';
import TemplatesList from '../../../components/documents/TemplatesList.vue';

useHead({
    title: 'Documents: Templates',
});
definePageMeta({
    requiresAuth: true,
    permission: 'DocStoreService.FindDocuments',
});
</script>

<template>
    <ContentWrapper>
        <TemplatesList />
    </ContentWrapper>
</template>

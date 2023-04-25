<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import CategoriesList from '~/components/documents/categories/CategoriesList.vue';

useHead({
    title: 'pages.documents.categories.title',
});
definePageMeta({
    title: 'pages.documents.categories.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListDocumentCategories',
});
</script>

<template>
    <ContentWrapper>
        <CategoriesList />
    </ContentWrapper>
</template>

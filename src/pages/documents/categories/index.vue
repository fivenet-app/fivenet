<script lang="ts" setup>
import CategoriesList from '~/components/documents/categories/CategoriesList.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'pages.documents.categories.title',
});
definePageMeta({
    title: 'pages.documents.categories.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListCategories',
});
</script>

<template>
    <ContentWrapper>
        <CategoriesList />
    </ContentWrapper>
</template>

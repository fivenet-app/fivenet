<script lang="ts" setup>
import ContentWrapper from '../../../components/partials/ContentWrapper.vue';
import CategoriesList from '../../../components/documents/categories/CategoriesList.vue';

useHead({
    title: 'Categories',
});
definePageMeta({
    title: 'Categories',
    requiresAuth: true,
    permission: 'CompletorService.CompleteDocumentCategory',
});
</script>

<template>
    <ContentWrapper>
        <CategoriesList />
    </ContentWrapper>
</template>

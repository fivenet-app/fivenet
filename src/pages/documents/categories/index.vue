<script lang="ts" setup>
import List from '~/components/documents/categories/List.vue';
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
        <List />
    </ContentWrapper>
</template>

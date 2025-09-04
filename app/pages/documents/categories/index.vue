<script lang="ts" setup>
import List from '~/components/documents/categories/List.vue';

useHead({
    title: 'pages.documents.categories.title',
});

definePageMeta({
    title: 'pages.documents.categories.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListCategories',
});
</script>

<template>
    <List />
</template>

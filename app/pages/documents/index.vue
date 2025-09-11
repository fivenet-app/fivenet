<script lang="ts" setup>
import List from '~/components/documents/List.vue';

useHead({
    title: 'pages.documents.title',
});

definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListDocuments',
});
</script>

<template>
    <List />
</template>

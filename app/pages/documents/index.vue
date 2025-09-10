<script lang="ts" setup>
import DocumentList from '~/components/documents/DocumentList.vue';

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
    <DocumentList />
</template>

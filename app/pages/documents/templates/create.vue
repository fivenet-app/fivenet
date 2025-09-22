<script lang="ts" setup>
import Editor from '~/components/documents/templates/Editor.vue';

useHead({
    title: 'pages.documents.templates.create_template',
});

definePageMeta({
    title: 'pages.documents.templates.create_template',
    requiresAuth: true,
    permission: 'documents.DocumentsService/CreateTemplate',
});
</script>

<template>
    <Editor />
</template>

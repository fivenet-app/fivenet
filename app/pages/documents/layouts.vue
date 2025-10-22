<script setup lang="ts">
import EditorContainer from '~/components/fabriceditor/EditorContainer.vue';

useHead({
    title: 'Document Layouts',
});

definePageMeta({
    title: 'Document Layouts',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});
</script>

<template>
    <EditorContainer />
</template>

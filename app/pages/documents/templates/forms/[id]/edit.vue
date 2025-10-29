<script setup lang="ts">
import EditorContainer from '~/components/fabriceditor/EditorContainer.vue';

useHead({
    title: 'Document Edit Form',
});

definePageMeta({
    title: 'Document Edit Form',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});
</script>

<template>
    <EditorContainer />
</template>

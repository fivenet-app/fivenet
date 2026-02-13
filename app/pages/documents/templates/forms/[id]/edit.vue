<script setup lang="ts">
import EditorContainer from '~/components/fabriceditor/EditorContainer.vue';

useHead({
    title: 'Edit Document Form',
});

definePageMeta({
    title: 'Edit Document Form',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});

const svgData = ref<string>('');
</script>

<template>
    <EditorContainer v-model="svgData" />
</template>

<script setup lang="ts">
import EditorContainer from '~/components/fabriceditor/EditorContainer.vue';

useHead({
    title: 'Document Forms',
});

definePageMeta({
    title: 'Document Forms',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});

const svgData = ref<string>('');
</script>

<template>
    <EditorContainer v-model="svgData" />
</template>

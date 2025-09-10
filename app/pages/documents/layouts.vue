<script lang="ts" setup>
import LayoutEditor from '~/components/documents/layouteditor/LayoutEditor.vue';

useHead({
    title: 'Layout Editor',
});

definePageMeta({
    title: 'Layout Editor',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});
</script>

<template>
    <LayoutEditor />
</template>

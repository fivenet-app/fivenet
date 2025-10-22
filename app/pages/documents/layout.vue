<script setup lang="ts">
import CanvasRenderer from '~/components/fabriceditor/CanvasRenderer.vue';

useHead({
    title: 'Document Layout',
});

definePageMeta({
    title: 'Document Layout',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});

const { data: json } = await useFetch<string>('/test.json');
</script>

<template>
    <div>
        <CanvasRenderer v-if="json" :json="json" />
    </div>
</template>

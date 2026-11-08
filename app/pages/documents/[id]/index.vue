<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import View from '~/components/documents/View.vue';

useHead({
    title: 'pages.documents.id.title',
});

definePageMeta({
    title: 'pages.documents.id.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListDocuments',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('documents-id');
const documentId = computed(() => (typeof route.params.id === 'string' ? parseInt(route.params.id) : 0));
</script>

<template>
    <View :document-id="documentId" />
</template>

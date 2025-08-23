<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import DocumentView from '~/components/documents/DocumentView.vue';

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
</script>

<template>
    <UDashboardPanel>
        <DocumentView :document-id="parseInt(route.params.id)" />
    </UDashboardPanel>
</template>

<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import DocumentEditor from '~/components/documents/DocumentEditor.vue';

useHead({
    title: 'pages.documents.edit.title',
});
definePageMeta({
    title: 'pages.documents.edit.title',
    requiresAuth: true,
    permission: 'DocStoreService.CreateDocument',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-id-edit'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('documents-id-edit');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <DocumentEditor :document-id="parseInt(route.params.id)" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

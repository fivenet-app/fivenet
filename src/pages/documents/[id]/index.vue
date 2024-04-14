<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import DocumentView from '~/components/documents/DocumentView.vue';

useHead({
    title: 'pages.documents.id.title',
});
definePageMeta({
    title: 'pages.documents.id.title',
    requiresAuth: true,
    permission: 'DocStoreService.GetDocument',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-id');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <DocumentView :document-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

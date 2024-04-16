<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import TemplateView from '~/components/documents/templates/TemplateView.vue';

useHead({
    title: 'pages.documents.templates.view.title',
});
definePageMeta({
    title: 'pages.documents.templates.view.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListTemplates',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-templates-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-templates-id');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <TemplateView :template-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

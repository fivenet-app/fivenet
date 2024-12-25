<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import TemplateEditor from '~/components/documents/templates/TemplateEditor.vue';

useHead({
    title: 'pages.documents.templates.edit.title',
});
definePageMeta({
    title: 'pages.documents.templates.edit.title',
    requiresAuth: true,
    permission: 'DocStoreService.CreateTemplate',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-templates-edit-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('documents-templates-edit-id');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <TemplateEditor :template-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

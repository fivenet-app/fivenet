<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
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
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-templates-edit-id');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.documents.templates.edit.title')"> </UDashboardNavbar>

            <TemplateEditor :template-id="route.params.id as string" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

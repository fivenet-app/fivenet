<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import Editor from '~/components/documents/templates/Editor.vue';

useHead({
    title: 'pages.documents.templates.edit.title',
});

definePageMeta({
    title: 'pages.documents.templates.edit.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/CreateTemplate',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-templates-id-edit'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('documents-templates-id-edit');
</script>

<template>
    <Editor :template-id="parseInt(route.params.id)" />
</template>

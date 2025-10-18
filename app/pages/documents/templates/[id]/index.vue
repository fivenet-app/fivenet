<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import View from '~/components/documents/templates/View.vue';

useHead({
    title: 'pages.documents.templates.view.title',
});

definePageMeta({
    title: 'pages.documents.templates.view.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListTemplates',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-templates-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('documents-templates-id');
</script>

<template>
    <View :template-id="parseInt(route.params.id)" />
</template>

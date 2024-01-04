<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import TemplateView from '~/components/documents/templates/TemplateView.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

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
    <ContentWrapper>
        <TemplateView :template-id="route.params.id as string" />
    </ContentWrapper>
</template>

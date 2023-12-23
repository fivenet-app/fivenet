<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import View from '~/components/documents/templates/View.vue';
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
        <View :template-id="route.params.id as string" />
    </ContentWrapper>
</template>

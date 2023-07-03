<script lang="ts" setup>
import Editor from '~/components/documents/templates/Editor.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

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
    <ContentWrapper>
        <Editor :templateId="BigInt(route.params.id)" />
    </ContentWrapper>
</template>

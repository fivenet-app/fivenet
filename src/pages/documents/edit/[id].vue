<script lang="ts" setup>
import DocumentEditor from '~/components/documents/DocumentEditor.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'pages.documents.edit.title',
});
definePageMeta({
    title: 'pages.documents.edit.title',
    requiresAuth: true,
    permission: 'DocStoreService.CreateDocument',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-edit-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-edit-id');
</script>

<template>
    <ContentWrapper>
        <DocumentEditor :id="parseInt(route.params.id)" />
    </ContentWrapper>
</template>

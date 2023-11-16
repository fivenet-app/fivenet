<script lang="ts" setup>
import Editor from '~/components/documents/Editor.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import { type TypedRouteFromName } from '@typed-router';

useHead({
    title: 'pages.documents.edit.title',
});
definePageMeta({
    title: 'pages.documents.edit.title',
    requiresAuth: true,
    permission: 'DocStoreService.CreateDocument',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-id-edit'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-id-edit');
</script>

<template>
    <ContentWrapper>
        <Editor :id="route.params.id as string" />
    </ContentWrapper>
</template>

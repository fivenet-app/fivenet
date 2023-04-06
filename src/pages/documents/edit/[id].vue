<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import DocumentEditor from '~/components/documents/DocumentEditor.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'Edit Document',
});
definePageMeta({
    title: 'Edit Document',
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

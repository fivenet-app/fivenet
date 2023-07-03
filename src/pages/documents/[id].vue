<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import View from '~/components/documents/View.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'pages.documents.id.title',
});
definePageMeta({
    title: 'pages.documents.id.title',
    requiresAuth: true,
    permission: 'DocStoreService.GetDocument',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-id');
</script>

<template>
    <ContentWrapper>
        <View :documentId="BigInt(route.params.id)" />
        <ClipboardButton />
    </ContentWrapper>
</template>

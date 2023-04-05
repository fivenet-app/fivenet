<script lang="ts" setup>
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';
import ContentWrapper from '../../../components/partials/ContentWrapper.vue';

useHead({
    title: 'Category: View',
});
definePageMeta({
    title: 'Category: View',
    requiresAuth: true,
    permission: 'DocStoreService.CreateDocumentCategory',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-categories-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});
</script>
<template>
    <ContentWrapper>
        <h2 class="text-2xl text-white">Category: View</h2>
        <!-- TODO -->
    </ContentWrapper>
</template>

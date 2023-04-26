<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import TemplateView from '~/components/documents/templates/TemplateView.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

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
const id = ref(0);

onMounted(() => {
    id.value = parseInt(route.params.id);
});
</script>

<template>
    <ContentWrapper>
        <TemplateView v-if="id > 0" :templateId="id" />
    </ContentWrapper>
</template>

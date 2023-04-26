<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import TemplateEditor from '~/components/documents/templates/TemplateEditor.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'pages.documents.templates.edit.title',
});
definePageMeta({
    title: 'pages.documents.templates.edit.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListTemplates',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-templates-edit-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-templates-edit-id');
const id = ref(0);

onMounted(() => {
    id.value = parseInt(route.params.id);
});
</script>

<template>
    <ContentWrapper>
        <TemplateEditor v-if="id > 0" :templateId="id" />
    </ContentWrapper>
</template>

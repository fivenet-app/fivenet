<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import TemplateView from '~/components/documents/templates/TemplateView.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'Template: View',
});
definePageMeta({
    title: 'Template: View',
    requiresAuth: true,
    permission: 'DocStoreService.ListTemplates',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-templates-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-templates-id');
const templateId = ref(0);

onMounted(() => {
    templateId.value = parseInt(route.params.id);
});
</script>

<template>
    <ContentWrapper>
        <TemplateView v-if="templateId > 0" :templateId="templateId" />
    </ContentWrapper>
</template>

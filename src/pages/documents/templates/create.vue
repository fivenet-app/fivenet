<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import TemplateEditor from '~/components/documents/templates/TemplateEditor.vue';

useHead({
    title: 'pages.documents.templates.create_template',
});
definePageMeta({
    title: 'pages.documents.templates.create_template',
    requiresAuth: true,
    permission: 'DocStoreService.CreateTemplate',
});
</script>

<template>
    <ContentWrapper>
        <TemplateEditor />
    </ContentWrapper>
</template>

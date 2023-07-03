<script lang="ts" setup>
import Editor from '~/components/documents/templates/Editor.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

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
        <Editor />
    </ContentWrapper>
</template>

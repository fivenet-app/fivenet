<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import TemplateEditor from '~/components/documents/templates/TemplateEditor.vue';

useHead({
    title: 'Create Template',
});
definePageMeta({
    title: 'Create Template',
    requiresAuth: true,
    permission: 'DocStoreService.CreateTemplate',
});
</script>

<template>
    <ContentWrapper>
        <TemplateEditor />
    </ContentWrapper>
</template>

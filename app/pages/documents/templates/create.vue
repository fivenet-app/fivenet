<script lang="ts" setup>
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

        <UDashboardPanel grow>
            <TemplateEditor />
        </UDashboardPanel>

</template>

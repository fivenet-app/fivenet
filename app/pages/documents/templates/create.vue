<script lang="ts" setup>
import TemplateEditor from '~/components/documents/templates/TemplateEditor.vue';

useHead({
    title: 'pages.documents.templates.create_template',
});

definePageMeta({
    title: 'pages.documents.templates.create_template',
    requiresAuth: true,
    permission: 'documents.DocumentsService.CreateTemplate',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <TemplateEditor />
        </UDashboardPanel>
    </UDashboardPage>
</template>

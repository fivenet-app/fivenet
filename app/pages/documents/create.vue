<script lang="ts" setup>
import DocumentEditor from '~/components/documents/DocumentEditor.vue';

useHead({
    title: 'pages.documents.create.title',
});

definePageMeta({
    title: 'pages.documents.create.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService.CreateDocument',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <DocumentEditor />
        </UDashboardPanel>
    </UDashboardPage>
</template>

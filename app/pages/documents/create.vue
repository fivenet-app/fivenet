<script lang="ts" setup>
import DocumentEditor from '~/components/documents/DocumentEditor.vue';

useHead({
    title: 'pages.documents.create.title',
});

definePageMeta({
    title: 'pages.documents.create.title',
    requiresAuth: true,
    permission: 'DocStoreService.CreateDocument',
});
</script>

<template>

        <UDashboardPanel grow>
            <DocumentEditor />
        </UDashboardPanel>

</template>

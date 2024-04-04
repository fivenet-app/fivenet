<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
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
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.documents.create.title')"> </UDashboardNavbar>

            <DocumentEditor />

            <ClipboardButton />
        </UDashboardPanel>
    </UDashboardPage>
</template>

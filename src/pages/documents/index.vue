<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import DocumentList from '~/components/documents/DocumentList.vue';

useHead({
    title: 'pages.documents.title',
});
definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListDocuments',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #right>
                    <UButton :label="$t('common.create')" trailing-icon="i-mdi-plus" color="gray" to="/documents/create" />
                </template>
            </UDashboardNavbar>

            <DocumentList />
            <ClipboardButton />
        </UDashboardPanel>
    </UDashboardPage>
</template>

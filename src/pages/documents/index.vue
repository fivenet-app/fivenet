<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import DocumentList from '~/components/documents/DocumentList.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

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
    <ContentWrapper>
        <DocumentList />
        <ClipboardButton />
    </ContentWrapper>
</template>

<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import List from '~/components/documents/List.vue';
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
        <List />
        <ClipboardButton />
    </ContentWrapper>
</template>

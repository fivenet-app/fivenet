<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentsList from '../../components/documents/DocumentsList.vue';
import ClipboardButton from '../../components/clipboard/ClipboardButton.vue';

useHead({
    title: 'Documents',
});
definePageMeta({
    requiresAuth: true,
    permission: 'DocStoreService.FindDocuments',
});
</script>

<template>
    <ContentWrapper>
        <DocumentsList />
        <ClipboardButton />
    </ContentWrapper>
</template>

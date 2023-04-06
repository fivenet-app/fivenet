<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import DocumentEditor from '~/components/documents/DocumentEditor.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';

useHead({
    title: 'Create Document',
});
definePageMeta({
    title: 'Create Document',
    requiresAuth: true,
    permission: 'DocStoreService.CreateDocument',
});
</script>

<template>
    <ContentWrapper>
        <DocumentEditor />
        <ClipboardButton />
    </ContentWrapper>
</template>

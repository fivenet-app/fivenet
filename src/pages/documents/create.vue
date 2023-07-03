<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import Editor from '~/components/documents/Editor.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

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
    <ContentWrapper>
        <Editor />
        <ClipboardButton />
    </ContentWrapper>
</template>

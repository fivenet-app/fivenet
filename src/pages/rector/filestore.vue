<script lang="ts" setup>
import FileList from '~/components/rector/filestore/FileList.vue';

useHead({
    title: 'pages.rector.filestore.title',
});
definePageMeta({
    title: 'pages.rector.filestore.title',
    requiresAuth: true,
    permission: 'SuperUser',
    showQuickButtons: false,
});
</script>

<template>
    <div>
        <FileList />
    </div>
</template>

<script lang="ts" setup>
import FileList from '~/components/settings/filestore/FileList.vue';

useHead({
    title: 'pages.settings.filestore.title',
});

definePageMeta({
    title: 'pages.settings.filestore.title',
    requiresAuth: true,
    permission: 'Superuser/Superuser',
});
</script>

<template>
    <FileList />
</template>

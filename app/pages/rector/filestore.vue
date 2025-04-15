<script lang="ts" setup>
import FileList from '~/components/rector/filestore/FileList.vue';

useHead({
    title: 'pages.rector.filestore.title',
});

definePageMeta({
    title: 'pages.rector.filestore.title',
    requiresAuth: true,
    permission: 'SuperUser',
});
</script>

<template>

        <UDashboardPanel grow>
            <FileList />
        </UDashboardPanel>

</template>

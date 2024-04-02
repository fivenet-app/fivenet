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
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.rector.filestore.title')"> </UDashboardNavbar>

            <FileList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

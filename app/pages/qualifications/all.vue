<script lang="ts" setup>
import List from '~/components/qualifications/List.vue';
import ListPage from '~/components/qualifications/ListPage.vue';

useHead({
    title: 'pages.qualifications.title',
});

definePageMeta({
    title: 'pages.qualifications.title',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
});
</script>

<template>
    <ListPage>
        <List />
    </ListPage>
</template>

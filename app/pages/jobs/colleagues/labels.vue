<script lang="ts" setup>
import List from '~/components/jobs/colleagues/labels/List.vue';

useHead({
    title: 'pages.jobs.colleagues.labels.title',
});

definePageMeta({
    title: 'pages.jobs.colleagues.labels.title',
    requiresAuth: true,
    permission: 'jobs.ColleaguesService/ManageLabels',
});
</script>

<template>
    <List />
</template>

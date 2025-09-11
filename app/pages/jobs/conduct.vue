<script lang="ts" setup>
import List from '~/components/jobs/conduct/List.vue';

useHead({
    title: 'pages.jobs.conduct.title',
});

definePageMeta({
    title: 'pages.jobs.conduct.title',
    requiresAuth: true,
    permission: 'jobs.ConductService/ListConductEntries',
});
</script>

<template>
    <List />
</template>

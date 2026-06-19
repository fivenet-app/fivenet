<script lang="ts" setup>
import List from '~/components/jobs/groups/List.vue';

useHead({
    title: 'pages.jobs.groups.title',
});

definePageMeta({
    title: 'pages.jobs.groups.title',
    requiresAuth: true,
    permission: 'jobs.GroupsService/ListGroups',
});
</script>

<template>
    <List />
</template>

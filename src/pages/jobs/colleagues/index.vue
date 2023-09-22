<script lang="ts" setup>
import List from '~/components/jobs/colleagues/List.vue';

useHead({
    title: 'pages.jobs.colleagues.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'JobsService.ColleaguesList',
});
</script>

<template>
    <div>
        <List />
    </div>
</template>

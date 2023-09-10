<script lang="ts" setup>
import List from '~/components/jobs/colleagues/List.vue';

useHead({
    title: 'pages.jobs.colleagues.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'Jobs.ColleaguesList',
});
</script>

<template>
    <div>
        <List :current-job="true" />
    </div>
</template>

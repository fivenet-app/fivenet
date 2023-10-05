<script lang="ts" setup>
import List from '~/components/jobs/requests/List.vue';

useHead({
    title: 'pages.jobs.requests.title',
});
definePageMeta({
    title: 'pages.jobs.requests.title',
    requiresAuth: true,
    permission: 'JobsService.RequestsList',
});
</script>

<template>
    <div>
        <List />
    </div>
</template>

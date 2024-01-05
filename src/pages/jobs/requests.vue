<script lang="ts" setup>
import RequestsList from '~/components/jobs/requests/RequestsList.vue';

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
        <RequestsList />
    </div>
</template>

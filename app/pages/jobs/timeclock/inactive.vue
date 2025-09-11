<script lang="ts" setup>
import InactiveList from '~/components/jobs/timeclock/InactiveList.vue';

useHead({
    title: 'common.inactive_colleagues',
});

definePageMeta({
    title: 'common.inactive_colleagues',
    requiresAuth: true,
    permission: 'jobs.TimeclockService/ListInactiveEmployees',
});
</script>

<template>
    <InactiveList />
</template>

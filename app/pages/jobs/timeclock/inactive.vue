<script lang="ts" setup>
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';
import TimeclockInactiveList from '~/components/jobs/timeclock/TimeclockInactiveList.vue';

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
    <PagesJobsLayout>
        <template #default>
            <TimeclockInactiveList />
        </template>
    </PagesJobsLayout>
</template>

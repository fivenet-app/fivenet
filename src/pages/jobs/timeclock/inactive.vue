<script lang="ts" setup>
import TimeclockInactiveList from '~/components/jobs/timeclock/TimeclockInactiveList.vue';
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';

useHead({
    title: 'common.inactive_colleagues',
});
definePageMeta({
    title: 'common.inactive_colleagues',
    requiresAuth: true,
    permission: 'JobsTimeclockService.ListInactiveEmployees',
});
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <TimeclockInactiveList />
        </template>
    </PagesJobsLayout>
</template>

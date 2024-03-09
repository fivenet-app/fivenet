<script lang="ts" setup>
import TimeclockInactiveList from '~/components/jobs/timeclock/TimeclockInactiveList.vue';

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
    <div>
        <TimeclockInactiveList />
    </div>
</template>

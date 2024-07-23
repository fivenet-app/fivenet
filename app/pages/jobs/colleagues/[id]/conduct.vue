<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ConductList from '~/components/jobs/conduct/ConductList.vue';

useHead({
    title: 'pages.jobs.colleagues.single.conduct',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.conduct',
    requiresAuth: true,
    permission: 'JobsConductService.ListConductEntries',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-conduct'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-colleagues-id-conduct');
</script>

<template>
    <div>
        <ConductList :user-id="parseInt(route.params.id)" :hide-user-search="true" />
    </div>
</template>

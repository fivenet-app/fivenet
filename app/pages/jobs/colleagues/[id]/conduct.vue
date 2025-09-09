<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ConductList from '~/components/jobs/conduct/ConductList.vue';

useHead({
    title: 'pages.jobs.colleagues.id.conduct',
});

definePageMeta({
    title: 'pages.jobs.colleagues.id.conduct',
    requiresAuth: true,
    permission: 'jobs.ConductService/ListConductEntries',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-conduct'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('jobs-colleagues-id-conduct');
</script>

<template>
    <ConductList :user-id="parseInt(route.params.id as string)" :hide-user-search="true" />
</template>

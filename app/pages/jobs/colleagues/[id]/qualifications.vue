<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ResultList from '~/components/qualifications/result/ResultList.vue';

useHead({
    title: 'pages.jobs.colleagues.id.qualifications',
});

definePageMeta({
    title: 'pages.jobs.colleagues.id.qualifications',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-qualifications'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('jobs-colleagues-id-qualifications');
</script>

<template>
    <ResultList class="m-4" :user-id="parseInt(route.params.id as string)" />
</template>

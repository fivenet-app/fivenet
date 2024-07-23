<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import QualificationsResultsList from '~/components/qualifications/QualificationsResultsList.vue';

useHead({
    title: 'pages.jobs.colleagues.single.qualifications',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.qualifications',
    requiresAuth: true,
    permission: 'QualificationsService.ListQualifications',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-qualifications'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-colleagues-id-qualifications');
</script>

<template>
    <QualificationsResultsList :user-id="parseInt(route.params.id)" />
</template>

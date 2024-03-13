<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import QualificationView from '~/components/jobs/qualifications/QualificationView.vue';

useHead({
    title: 'pages.jobs.qualifications.single.title',
});
definePageMeta({
    title: 'pages.jobs.qualifications.single.title',
    requiresAuth: true,
    permission: 'JobsQualificationsService.GetQualification',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-qualifications-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-qualifications-id');
</script>

<template>
    <div>
        <QualificationView :id="route.params.id" />
    </div>
</template>

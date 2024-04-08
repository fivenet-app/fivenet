<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import QualificationView from '~/components/jobs/qualifications/QualificationView.vue';

useHead({
    title: 'pages.qualifications.single.title',
});
definePageMeta({
    title: 'pages.qualifications.single.title',
    requiresAuth: true,
    permission: 'QualificationsService.GetQualification',
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
        <QualificationView :qualification-id="route.params.id" />
    </div>
</template>

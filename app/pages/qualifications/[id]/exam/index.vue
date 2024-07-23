<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import ExamView from '~/components/qualifications/exam/ExamView.vue';

useHead({
    title: 'pages.qualifications.single.exam.title',
});
definePageMeta({
    title: 'pages.qualifications.single.exam.title',
    requiresAuth: true,
    permission: 'QualificationsService.GetQualification',
    validate: async (route) => {
        route = route as TypedRouteFromName<'qualifications-id-exam'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('qualifications-id');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <ExamView :qualification-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

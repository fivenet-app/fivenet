<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import QualificationEditor from '~/components/jobs/qualifications/QualificationEditor.vue';

useHead({
    title: 'pages.qualifications.edit.title',
});
definePageMeta({
    title: 'pages.qualifications.edit.title',
    requiresAuth: true,
    permission: 'QualificationsService.UpdateQualification',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-qualifications-id-edit'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-qualifications-id-edit');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.qualifications.edit.title')"> </UDashboardNavbar>

            <QualificationEditor :qualification-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

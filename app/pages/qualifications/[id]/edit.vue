<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import QualificationEditor from '~/components/qualifications/QualificationEditor.vue';

useHead({
    title: 'pages.qualifications.edit.title',
});
definePageMeta({
    title: 'pages.qualifications.edit.title',
    requiresAuth: true,
    permission: 'QualificationsService.UpdateQualification',
    validate: async (route) => {
        route = route as TypedRouteFromName<'qualifications-id-edit'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('qualifications-id-edit');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <QualificationEditor :qualification-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

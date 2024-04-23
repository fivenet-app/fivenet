<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
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
        return /^\d+$/.test(route.params.id);
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

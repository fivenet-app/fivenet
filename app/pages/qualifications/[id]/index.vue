<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import QualificationView from '~/components/qualifications/QualificationView.vue';

useHead({
    title: 'pages.qualifications.single.title',
});
definePageMeta({
    title: 'pages.qualifications.single.title',
    requiresAuth: true,
    permission: 'QualificationsService.ListQualifications',
    validate: async (route) => {
        route = route as TypedRouteFromName<'qualifications-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return idParamRegex.test(route.params.id as string);
    },
});

const route = useRoute('qualifications-id');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <QualificationView :qualification-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

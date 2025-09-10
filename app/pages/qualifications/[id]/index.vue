<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import QualificationView from '~/components/qualifications/QualificationView.vue';

useHead({
    title: 'pages.qualifications.id.title',
});

definePageMeta({
    title: 'pages.qualifications.id.title',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
    validate: async (route) => {
        route = route as TypedRouteFromName<'qualifications-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('qualifications-id');
</script>

<template>
    <QualificationView :qualification-id="parseInt(route.params.id)" />
</template>

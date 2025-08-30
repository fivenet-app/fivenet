<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ExamView from '~/components/qualifications/exam/ExamView.vue';

useHead({
    title: 'pages.qualifications.single.exam.title',
});

definePageMeta({
    title: 'pages.qualifications.single.exam.title',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
    validate: async (route) => {
        route = route as TypedRouteFromName<'qualifications-id-exam'>;
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
    <ExamView :qualification-id="parseInt(route.params.id)" />
</template>

<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import ExamEditor from '~/components/qualifications/exam/ExamEditor.vue';

useHead({
    title: 'pages.qualifications.single.exam.edit.title',
});
definePageMeta({
    title: 'pages.qualifications.single.exam.edit.title',
    requiresAuth: true,
    permission: 'QualificationsService.UpdateQualification',
    validate: async (route) => {
        route = route as TypedRouteFromName<'qualifications-id-exam-edit'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('qualifications-id-exam-edit');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <ExamEditor :qualification-id="route.params.id" />
        </UDashboardPanel>
    </UDashboardPage>
</template>

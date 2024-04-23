<script lang="ts" setup>
import QualificationEditor from '~/components/qualifications/QualificationEditor.vue';

useHead({
    title: 'pages.qualifications.single.title',
});
definePageMeta({
    title: 'pages.qualifications.single.title',
    requiresAuth: true,
    permission: 'QualificationsService.CreateQualification',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <QualificationEditor />
        </UDashboardPanel>
    </UDashboardPage>
</template>

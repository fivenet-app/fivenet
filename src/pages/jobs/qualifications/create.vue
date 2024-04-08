<script lang="ts" setup>
import QualificationEditor from '~/components/jobs/qualifications/QualificationEditor.vue';

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
            <UDashboardNavbar :title="$t('pages.qualifications.single.title')"> </UDashboardNavbar>

            <QualificationEditor />
        </UDashboardPanel>
    </UDashboardPage>
</template>

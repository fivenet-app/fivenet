<script lang="ts" setup>
import QualificationEditor from '~/components/qualifications/QualificationEditor.vue';

useHead({
    title: 'pages.qualifications.create',
});

definePageMeta({
    title: 'pages.qualifications.create',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService.CreateQualification',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <QualificationEditor />
        </UDashboardPanel>
    </UDashboardPage>
</template>

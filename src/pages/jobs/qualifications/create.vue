<script lang="ts" setup>
import QualificationEditor from '~/components/jobs/qualifications/QualificationEditor.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'pages.jobs.qualifications.single.title',
});
definePageMeta({
    title: 'pages.jobs.qualifications.single.title',
    requiresAuth: true,
    permission: 'QualificationsService.CreateQualification',
});
</script>

<template>
    <ContentWrapper>
        <QualificationEditor />
    </ContentWrapper>
</template>

<script lang="ts" setup>
import ConductList from '~/components/jobs/conduct/ConductList.vue';
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';

useHead({
    title: 'pages.jobs.conduct.title',
});

definePageMeta({
    title: 'pages.jobs.conduct.title',
    requiresAuth: true,
    permission: 'jobs.ConductService.ListConductEntries',
});
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <ConductList />
        </template>
    </PagesJobsLayout>
</template>

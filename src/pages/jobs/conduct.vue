<script lang="ts" setup>
import ConductList from '~/components/jobs/conduct/ConductList.vue';

useHead({
    title: 'pages.jobs.conduct.title',
});
definePageMeta({
    title: 'pages.jobs.conduct.title',
    requiresAuth: true,
    permission: 'JobsConductService.ListConductEntries',
});
</script>

<template>
    <div>
        <ConductList />
    </div>
</template>

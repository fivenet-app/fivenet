<script lang="ts" setup>
useHead({
    title: 'pages.jobs.title',
});

definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: [
        'jobs.ConductService/CreateConductEntry',
        'jobs.ConductService/DeleteConductEntry',
        'jobs.ConductService/ListConductEntries',
        'jobs.ConductService/UpdateConductEntry',
        'jobs.ConductService/UploadFile',
        'jobs.JobsService/GetColleague',
        'jobs.JobsService/ListColleagueActivity',
        'jobs.JobsService/ListColleagues',
        'jobs.StatsService/GetStats',
        'jobs.JobsService/ManageLabels',
        'jobs.JobsService/SetColleagueProps',
        'jobs.JobsService/SetMOTD',
        'jobs.TimeclockService/ListInactiveEmployees',
        'jobs.TimeclockService/ListTimeclock',
    ],
    redirect: { name: 'jobs-overview' },
});
</script>

<template>
    <div></div>
</template>

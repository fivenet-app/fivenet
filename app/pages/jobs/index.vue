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
        'jobs.ColleaguesService/GetColleague',
        'jobs.ColleaguesService/ListColleagueActivity',
        'jobs.ColleaguesService/ListColleagues',
        'jobs.StatsService/GetStats',
        'jobs.ColleaguesService/ManageLabels',
        'jobs.ColleaguesService/SetColleagueProps',
        'jobs.JobsService/SetMOTD',
        'jobs.TimeclockService/ListInactiveEmployees',
        'jobs.TimeclockService/ListTimeclock',
    ],
    redirect: '/jobs/overview',
});
</script>

<template>
    <div></div>
</template>

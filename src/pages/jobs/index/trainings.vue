<script lang="ts" setup>
useHead({
    title: 'pages.jobs.trainings.title',
});
definePageMeta({
    title: 'pages.jobs.trainings.title',
    requiresAuth: true,
    permission: 'Jobs.View',
});
</script>

<template>
    <h1 class="text-xl text-white">Trainings</h1>
</template>

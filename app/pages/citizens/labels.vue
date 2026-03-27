<script lang="ts" setup>
import List from '~/components/citizens/labels/List.vue';

useHead({
    title: 'pages.citizens.labels.title',
});

definePageMeta({
    title: 'pages.citizens.labels.title',
    requiresAuth: true,
    permission: ['completor.CompletorService/CompleteCitizenLabels', 'citizens.LabelsService/CreateOrUpdateLabel'],
});
</script>

<template>
    <List />
</template>

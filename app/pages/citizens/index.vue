<script lang="ts" setup>
import CitizenList from '~/components/citizens/CitizenList.vue';

useHead({
    title: 'pages.citizens.title',
});

definePageMeta({
    title: 'pages.citizens.title',
    requiresAuth: true,
    permission: 'citizens.CitizensService/ListCitizens',
});
</script>

<template>
    <CitizenList />
</template>

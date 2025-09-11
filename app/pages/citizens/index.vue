<script lang="ts" setup>
import List from '~/components/citizens/List.vue';

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
    <List />
</template>

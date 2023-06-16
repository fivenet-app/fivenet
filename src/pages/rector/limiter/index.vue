<script lang="ts" setup>
import RolesList from '~/components/rector/RolesList.vue';

useHead({
    title: 'pages.rector.limiter.title',
});
definePageMeta({
    title: 'pages.rector.limiter.title',
    requiresAuth: true,
    permission: 'SuperUser',
    showQuickButtons: false,
});
</script>
<template>
    <RolesList :to="{ name: 'rector-limiter-id', params: { id: '0' } }" />
</template>

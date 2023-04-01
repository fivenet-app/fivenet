<script lang="ts" setup>
import RolesList from '~~/src/components/rector/RolesList.vue';

useHead({
    title: 'Roles',
});
definePageMeta({
    title: 'Roles',
    requiresAuth: true,
    permission: 'RectorService.GetRoles',
});
</script>

<template>
    <RolesList />
</template>

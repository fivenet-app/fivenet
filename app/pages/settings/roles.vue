<script lang="ts" setup>
import RoleList from '~/components/settings/roles/RoleList.vue';

useHead({
    title: 'pages.settings.roles.title',
});

definePageMeta({
    title: 'pages.settings.roles.title',
    requiresAuth: true,
    permission: 'settings.SettingsService/GetRoles',
});
</script>

<template>
    <RoleList />
</template>

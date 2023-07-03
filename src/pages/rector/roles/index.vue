<script lang="ts" setup>
import RolesList from '~/components/rector/roles/RolesList.vue';

useHead({
    title: 'pages.rector.roles.title',
});
definePageMeta({
    title: 'pages.rector.roles.title',
    requiresAuth: true,
    permission: 'RectorService.GetRoles',
    showQuickButtons: false,
});
</script>

<template>
    <div class="w-full">
        <RolesList />
    </div>
</template>

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
    <div class="w-full">
        <RolesList />
    </div>
</template>

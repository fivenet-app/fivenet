<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import RoleView from '~/components/rector/roles/RoleView.vue';

useHead({
    title: 'pages.rector.roles.id.title',
});
definePageMeta({
    title: 'pages.rector.roles.id.title',
    requiresAuth: true,
    permission: ['RectorService.GetRoles', 'RectorService.SetJobProps'],
    showQuickButtons: false,
    validate: async (route) => {
        route = route as TypedRouteFromName<'rector-roles-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('rector-roles-id');
</script>

<template>
    <div class="w-full">
        <RoleView :role-id="route.params.id as string" />
    </div>
</template>

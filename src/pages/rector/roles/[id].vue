<script lang="ts" setup>
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';
import RoleView from '~~/src/components/rector/RoleView.vue';

useHead({
    title: 'Role: Info',
});
definePageMeta({
    title: 'Role: Info',
    requiresAuth: true,
    permission: 'RectorService.GetRoles',
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
        <RoleView :role-id="parseInt(route.params.id)" />
    </div>
</template>

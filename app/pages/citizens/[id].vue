<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import CitizenInfo from '~/components/citizens/info/CitizenInfo.vue';

useHead({
    title: 'pages.citizens.id.title',
});
definePageMeta({
    title: 'pages.citizens.id.title',
    requiresAuth: true,
    permission: 'CitizenStoreService.GetUser',
    validate: async (route) => {
        route = route as TypedRouteFromName<'citizens-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('citizens-id');
</script>

<template>
    <CitizenInfo :user-id="parseInt(route.params.id)" />
</template>

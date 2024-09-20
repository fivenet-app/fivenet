<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
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
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return idParamRegex.test(route.params.id as string);
    },
});

const route = useRoute('citizens-id');
</script>

<template>
    <CitizenInfo :user-id="parseInt(route.params.id as string)" />
</template>

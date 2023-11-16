<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import Info from '~/components/citizens/info/Info.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'pages.citizens.id.title',
});
definePageMeta({
    title: 'pages.citizens.id.title',
    requiresAuth: true,
    permission: 'CitizenStoreService.ListCitizens',
    validate: async (route) => {
        route = route as TypedRouteFromName<'citizens-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('citizens-id');
</script>

<template>
    <ContentWrapper>
        <Info :id="route.params.id" />
    </ContentWrapper>
</template>

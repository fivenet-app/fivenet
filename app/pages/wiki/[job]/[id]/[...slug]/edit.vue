<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import PageEditor from '~/components/wiki/PageEditor.vue';

useHead({
    title: 'common.wiki',
});

definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
    permission: 'wiki.WikiService/ListPages',
    validate: async (route) => {
        route = route as TypedRouteFromName<'wiki-job-id-slug-edit'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('wiki-job-id-slug-edit');
</script>

<template>
    <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-r lg:border-b-0 dark:border-gray-800">
        <PageEditor :page-id="parseInt(route.params.id)" />
    </UDashboardPanel>
</template>

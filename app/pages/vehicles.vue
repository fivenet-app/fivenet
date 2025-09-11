<script lang="ts" setup>
import List from '~/components/vehicles/List.vue';

useHead({
    title: 'pages.vehicles.title',
});

definePageMeta({
    title: 'pages.vehicles.title',
    requiresAuth: true,
    permission: 'vehicles.VehiclesService/ListVehicles',
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.vehicles.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <List />
        </template>
    </UDashboardPanel>
</template>

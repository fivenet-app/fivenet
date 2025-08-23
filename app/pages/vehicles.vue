<script lang="ts" setup>
import VehicleList from '~/components/vehicles/VehicleList.vue';

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
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.vehicles.title')" />
        </template>

        <template #body>
            <VehicleList />
        </template>
    </UDashboardPanel>
</template>

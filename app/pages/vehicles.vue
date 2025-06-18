<script lang="ts" setup>
import VehiclesList from '~/components/vehicles/VehiclesList.vue';

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
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.vehicles.title')" />

            <VehiclesList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import VehiclesList from '~/components/vehicles/VehiclesList.vue';

useHead({
    title: 'pages.vehicles.title',
});
definePageMeta({
    title: 'pages.vehicles.title',
    requiresAuth: true,
    permission: 'DMVService.ListVehicles',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.vehicles.title')"> </UDashboardNavbar>

            <VehiclesList />
            <ClipboardButton />
        </UDashboardPanel>
    </UDashboardPage>
</template>

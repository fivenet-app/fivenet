<script lang="ts" setup>
import ClipboardButton from '../../components/clipboard/ClipboardButton.vue';
import VehiclesList from '../../components/vehicles/VehiclesList.vue';

useHead({
    title: 'Vehicles',
});
definePageMeta({
    title: 'Vehicles',
    requiresAuth: true,
    permission: 'DMVService.FindVehicles',
});
</script>

<template>
    <div class="w-full">
        <VehiclesList />
        <ClipboardButton />
    </div>
</template>

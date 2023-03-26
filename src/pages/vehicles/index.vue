<script lang="ts" setup>
import ClipboardButton from '../../components/clipboard/ClipboardButton.vue';
import VehiclesList from '../../components/vehicles/VehiclesList.vue';
</script>

<route lang="json">
{
    "name": "Vehicles",
    "meta": {
        "requiresAuth": true,
        "permission": "DMVService.FindVehicles"
    }
}
</route>

<template>
    <div class="w-full">
        <VehiclesList />
        <ClipboardButton />
    </div>
</template>

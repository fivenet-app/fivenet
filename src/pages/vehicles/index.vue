<script setup lang="ts">
import Footer from '../../components/partials/Footer.vue';
import VehiclesList from '../../components/vehicles/VehiclesList.vue';
</script>

<route lang="json">
{
    "name": "Vehicles",
    "meta": {
        "requiresAuth": true,
        "permission": "DMVService.FindVehicles",
        "breadCrumbs": null
    }
}
</route>

<template>
    <div class="w-full">
        <VehiclesList />
    </div>
    <Footer />
</template>

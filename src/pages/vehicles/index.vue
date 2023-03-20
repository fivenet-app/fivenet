<script setup lang="ts">
import Footer from '../../components/partials/Footer.vue';
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
        VEHICLES
    </div>
    <Footer />
</template>

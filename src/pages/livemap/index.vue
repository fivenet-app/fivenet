<script setup lang="ts">
import Footer from '../../components/partials/Footer.vue';
import Livemap from '../../components/Livemap.vue';
</script>

<route lang="json">
{
    "name": "Livemap",
    "meta": {
        "requiresAuth": true,
        "permission": "LivemapperService.Stream"
    }
}
</route>

<template>
    <div class="w-full">
        <Livemap />
    </div>
    <Footer />
</template>

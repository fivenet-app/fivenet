<script lang="ts" setup>
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

<style>
#map {
    height: 94vh;
}
</style>

<template>
    <div class="w-full">
        <Livemap />
    </div>
</template>

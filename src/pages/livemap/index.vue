<script lang="ts" setup>
import Livemap from '../../components/Livemap.vue';

useHead({
    title: 'Livemap',
});
definePageMeta({
    requiresAuth: true,
    permission: 'LivemapperService.Stream',
});
</script>

<style>
#map {
    height: 94.75vh;
}
</style>

<template>
    <div class="w-full">
        <Livemap />
    </div>
</template>

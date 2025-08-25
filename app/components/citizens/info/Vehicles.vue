<script lang="ts" setup>
import VehicleList from '~/components/vehicles/VehicleList.vue';

defineProps<{
    userId: number;
}>();
</script>

<template>
    <VehicleList :user-id="userId" :hide-owner="true" :hide-citizen-link="true" />
</template>

<script lang="ts" setup>
import VehiclesList from '~/components/vehicles/VehiclesList.vue';

defineProps<{
    userId: number;
}>();
</script>

<template>
    <VehiclesList :user-id="userId" :hide-owner="true" :hide-citizen-link="true" />
</template>

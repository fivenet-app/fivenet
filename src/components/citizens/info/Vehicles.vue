<script lang="ts" setup>
import List from '~/components/vehicles/VehiclesList.vue';

defineProps<{
    userId: number;
}>();
</script>

<template>
    <List :user-id="userId" :hide-owner="true" :hide-citizen-link="true" />
</template>

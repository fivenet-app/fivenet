<script lang="ts" setup>
import List from '~/components/vehicles/List.vue';

defineProps<{
    userId: number;
}>();
</script>

<template>
    <List :user-id="userId" hide-owner hide-citizen-link />
</template>

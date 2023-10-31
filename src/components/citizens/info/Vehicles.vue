<script lang="ts" setup>
import List from '~/components/vehicles/List.vue';

defineProps<{
    userId: number;
}>();
</script>

<template>
    <List :user-id="userId" :hide-owner="true" :hide-citizen-link="true" />
</template>

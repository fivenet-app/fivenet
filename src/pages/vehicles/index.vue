<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import List from '~/components/vehicles/List.vue';

useHead({
    title: 'pages.vehicles.title',
});
definePageMeta({
    title: 'pages.vehicles.title',
    requiresAuth: true,
    permission: 'DMVService.ListVehicles',
});
</script>

<template>
    <div class="w-full">
        <List />
        <ClipboardButton />
    </div>
</template>

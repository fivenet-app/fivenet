<script lang="ts" setup>
import UnitsList from '~/components/rector/UnitsList.vue';

useHead({
    title: 'common.units',
});
definePageMeta({
    title: 'common.units',
    requiresAuth: true,
    permission: 'UnitService.ListUnits',
    showQuickButtons: false,
});
</script>

<template>
    <div class="w-full">
        <UnitsList />
    </div>
</template>

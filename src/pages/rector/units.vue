<script lang="ts" setup>
import List from '~/components/rector/units/List.vue';

useHead({
    title: 'common.units',
});
definePageMeta({
    title: 'common.units',
    requiresAuth: true,
    permission: 'CentrumService.ListUnits',
    showQuickButtons: false,
});
</script>

<template>
    <div class="w-full">
        <List />
    </div>
</template>

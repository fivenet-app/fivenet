<script lang="ts" setup>
import List from '~/components/centrum/settings/List.vue';

useHead({
    title: 'common.units',
});
definePageMeta({
    title: 'common.units',
    requiresAuth: true,
    permission: 'CentrumService.CreateOrUpdateUnit',
    showQuickButtons: false,
});
</script>

<template>
    <div class="w-full">
        <List />
    </div>
</template>

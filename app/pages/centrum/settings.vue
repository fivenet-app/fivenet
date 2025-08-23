<script lang="ts" setup>
import Settings from '~/components/centrum/settings/Settings.vue';

useHead({
    title: 'components.centrum.settings.title',
});

definePageMeta({
    title: 'components.centrum.settings.title',
    requiresAuth: true,
    permission: 'centrum.CentrumService/TakeControl',
});
</script>

<template>
    <UDashboardPanel>
        <Settings />
    </UDashboardPanel>
</template>

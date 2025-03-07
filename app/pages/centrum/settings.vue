<script lang="ts" setup>
import CentrumSettings from '~/components/centrum/settings/CentrumSettings.vue';

useHead({
    title: 'components.centrum.settings.title',
});

definePageMeta({
    title: 'components.centrum.settings.title',
    requiresAuth: true,
    permission: 'CentrumService.TakeControl',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <CentrumSettings />
        </UDashboardPanel>
    </UDashboardPage>
</template>

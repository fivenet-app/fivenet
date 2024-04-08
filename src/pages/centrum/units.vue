<script lang="ts" setup>
import UnitsList from '~/components/centrum/settings/UnitsList.vue';

useHead({
    title: 'common.unit',
});
definePageMeta({
    title: 'common.unit',
    requiresAuth: true,
    permission: 'CentrumService.CreateOrUpdateUnit',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UnitsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

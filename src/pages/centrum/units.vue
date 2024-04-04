<script lang="ts" setup>
import UnitsList from '~/components/centrum/settings/UnitsList.vue';

useHead({
    title: 'common.units',
});
definePageMeta({
    title: 'common.units',
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

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
            <UDashboardNavbar :title="$t('common.units')"> </UDashboardNavbar>

            <UnitsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

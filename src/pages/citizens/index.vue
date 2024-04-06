<script lang="ts" setup>
import CitizensList from '~/components/citizens/CitizensList.vue';

useHead({
    title: 'pages.citizens.title',
});
definePageMeta({
    title: 'pages.citizens.title',
    requiresAuth: true,
    permission: 'CitizenStoreService.ListCitizens',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.citizens.title')"> </UDashboardNavbar>

            <CitizensList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

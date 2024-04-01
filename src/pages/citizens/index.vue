<script lang="ts" setup>
import CitizensList from '~/components/citizens/CitizensList.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';

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
            <ClipboardButton />
        </UDashboardPanel>
    </UDashboardPage>
</template>

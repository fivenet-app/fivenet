<script lang="ts" setup>
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import NotificationList from '~/components/notifications/NotificationList.vue';

useHead({
    title: 'components.notifications.title',
});

definePageMeta({
    title: 'components.notifications.title',
    requiresAuth: true,
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('components.notifications.title')">
                <template #right>
                    <DNBToggle />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <NotificationList />
        </template>
    </UDashboardPanel>
</template>

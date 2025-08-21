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
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('components.notifications.title')">
                <template #right>
                    <DNBToggle />
                </template>
            </UDashboardNavbar>

            <NotificationList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

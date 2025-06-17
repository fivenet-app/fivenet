<script lang="ts" setup>
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import NotificationsList from '~/components/notifications/NotificationsList.vue';

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

            <NotificationsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

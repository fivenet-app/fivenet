<script lang="ts" setup>
import NotificationsList from '~/components/partials/notification/NotificationsList.vue';

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
            <UDashboardNavbar :title="$t('components.notifications.title')" />

            <NotificationsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>

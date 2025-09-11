<script lang="ts" setup>
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import List from '~/components/notifications/List.vue';

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
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <DNBToggle />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <List />
        </template>
    </UDashboardPanel>
</template>

<script setup lang="ts">
import NotificationsList from './partials/notification/NotificationsList.vue';

const { isNotificationsSlideoverOpen } = useDashboard();
</script>

<template>
    <UDashboardSlideover v-model="isNotificationsSlideoverOpen" title="Notifications">
        <template #default>
            <NotificationsList @clicked="isNotificationsSlideoverOpen = false" />
        </template>

        <template #footer>
            <UButtonGroup class="flex w-full">
                <UButton
                    block
                    class="flex-1"
                    @click="
                        navigateTo({ name: 'notifications' });
                        isNotificationsSlideoverOpen = false;
                    "
                >
                    {{ $t('components.partials.sidebar_notifications') }}
                </UButton>
                <UButton block class="flex-1" color="gray" @click="isNotificationsSlideoverOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
        </template>
    </UDashboardSlideover>
</template>

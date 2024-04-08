<script setup lang="ts">
import NotificationsList from './partials/notification/NotificationsList.vue';

const { isNotificationsSlideoverOpen } = useDashboard();
</script>

<template>
    <UDashboardSlideover
        v-model="isNotificationsSlideoverOpen"
        :ui="{ width: 'w-screen max-w-4xl' }"
        :title="$t('common.notification', 2)"
    >
        <template #default>
            <NotificationsList @clicked="isNotificationsSlideoverOpen = false" />
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
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

                <UButton color="black" block class="flex-1" @click="isNotificationsSlideoverOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
        </template>
    </UDashboardSlideover>
</template>

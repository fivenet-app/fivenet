<script setup lang="ts">
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import NotificationList from '~/components/notifications/NotificationList.vue';

const { isNotificationSlideoverOpen } = useDashboard();
</script>

<template>
    <USlideover v-model:open="isNotificationSlideoverOpen" :title="$t('common.notification', 2)">
        <template #actions>
            <DNBToggle />
        </template>

        <template #body>
            <NotificationList @clicked="isNotificationSlideoverOpen = false" />
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    @click="
                        navigateTo({ name: 'notifications' });
                        isNotificationSlideoverOpen = false;
                    "
                >
                    {{ $t('components.partials.sidebar_notifications') }}
                </UButton>

                <UButton class="flex-1" color="neutral" block @click="isNotificationSlideoverOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
        </template>
    </USlideover>
</template>

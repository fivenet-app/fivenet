<script setup lang="ts">
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import NotificationList from '~/components/notifications/NotificationList.vue';

const { isNotificationSlideoverOpen } = useDashboard();
</script>

<template>
    <USlideover
        v-model:open="isNotificationSlideoverOpen"
        :title="$t('common.notification', 2)"
        :ui="{ body: 'flex flex-col p-1 sm:p-1' }"
    >
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
                    :label="$t('components.partials.sidebar_notifications')"
                    @click="
                        navigateTo({ name: 'notifications' });
                        isNotificationSlideoverOpen = false;
                    "
                />

                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :label="$t('common.close', 1)"
                    @click="isNotificationSlideoverOpen = false"
                />
            </UButtonGroup>
        </template>
    </USlideover>
</template>

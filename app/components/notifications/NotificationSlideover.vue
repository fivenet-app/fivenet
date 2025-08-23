<script setup lang="ts">
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import NotificationList from '~/components/notifications/NotificationList.vue';

const { isNotificationSlideoverOpen } = useDashboard();
</script>

<template>
    <USlideover v-model="isNotificationSlideoverOpen">
        <template #title>
            <div class="flex items-center justify-between">
                <h3 class="inline-flex gap-2 text-2xl font-semibold leading-6">
                    {{ $t('common.notification', 2) }}
                </h3>

                <DNBToggle />
            </div>
        </template>

        <template #body>
            <div class="flex flex-1 flex-col">
                <NotificationList @clicked="isNotificationSlideoverOpen = false" />
            </div>
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

<script setup lang="ts">
import { formatTimeAgo } from '@vueuse/core';
import NotificationsList from './partials/notification/NotificationsList.vue';

const { isNotificationsSlideoverOpen } = useDashboard();
</script>

<template>
    <UDashboardSlideover v-model="isNotificationsSlideoverOpen" title="Notifications">
        <template #default>
            <NuxtLink
                v-for="notification in []"
                :key="notification.id"
                to="/notifications"
                class="p-3 rounded-md hover:bg-gray-50 dark:hover:bg-gray-800/50 cursor-pointer flex items-center gap-3 relative"
            >
                <UChip color="red" :show="!!notification.unread" inset>
                    <UAvatar v-bind="notification.sender.avatar" :alt="notification.sender.name" size="md" />
                </UChip>

                <div class="text-sm flex-1">
                    <p class="flex items-center justify-between">
                        <span class="text-gray-900 dark:text-white font-medium">{{ notification.sender.name }}</span>

                        <time
                            :datetime="notification.date"
                            class="text-gray-500 dark:text-gray-400 text-xs"
                            v-text="formatTimeAgo(new Date(notification.date))"
                        />
                    </p>
                    <p class="text-gray-500 dark:text-gray-400">
                        {{ notification.body }}
                    </p>
                </div>
            </NuxtLink>
            <NotificationsList @clicked="open = false" />
        </template>

        <template #footer>
            <span class="isolate inline-flex w-full rounded-md pr-4 shadow-sm">
                <button
                    type="button"
                    class="relative inline-flex w-full items-center rounded-l-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400"
                    @click="
                        navigateTo({ name: 'notifications' });
                        isNotificationsSlideoverOpen = false;
                    "
                >
                    {{ $t('components.partials.sidebar_notifications') }}
                </button>
                <button
                    type="button"
                    class="relative -ml-px inline-flex w-full items-center rounded-r-md bg-neutral-50 px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-200 hover:text-gray-900"
                    @click="isNotificationsSlideoverOpen = false"
                >
                    {{ $t('common.close', 1) }}
                </button>
            </span>
        </template>
    </UDashboardSlideover>
</template>

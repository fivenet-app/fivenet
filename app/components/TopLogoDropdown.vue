<script lang="ts" setup>
defineProps<{
    collapsed?: boolean;
    hideNotifications?: boolean;
}>();

const { isNotificationSlideoverOpen } = useDashboard();

const notifications = useNotificationsStore();
const { notificationCount, doNotDisturb } = storeToRefs(notifications);

const newNotification = ref<boolean>(false);

const { start } = useTimeoutFn(() => (newNotification.value = false), 1000, {
    immediate: false,
});

const currentCount = ref(0);
watch(notificationCount, () => {
    if (notificationCount.value > currentCount.value) {
        newNotification.value = true;
        start();
    }

    currentCount.value = notificationCount.value;
});

const version = APP_VERSION;
</script>

<template>
    <UTooltip :text="version">
        <UButton
            class="w-full"
            color="neutral"
            variant="ghost"
            :avatar="{
                src: '/images/logo.webp',
                alt: 'FiveNet',
            }"
            :ui="{ base: collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)', leadingAvatar: 'rounded-none' }"
        >
            <span v-if="!collapsed" class="truncate font-semibold text-highlighted">FiveNet</span>
        </UButton>
    </UTooltip>

    <UTooltip v-if="!collapsed && !hideNotifications" :text="$t('components.partials.sidebar_notifications')" :kbds="['B']">
        <UChip
            :show="notificationCount > 0"
            color="error"
            inset
            :text="notificationCount <= 9 ? notificationCount : '9+'"
            size="xl"
        >
            <UButton
                color="neutral"
                variant="ghost"
                square
                :icon="
                    doNotDisturb
                        ? 'i-mdi-bell-off-outline'
                        : notificationCount === 0
                          ? 'i-mdi-notifications-none'
                          : 'i-mdi-notifications'
                "
                @click="() => (isNotificationSlideoverOpen = true)"
            />
        </UChip>
    </UTooltip>
</template>

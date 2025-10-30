<script lang="ts" setup>
defineProps<{
    collapsed?: boolean;
}>();

const { isNotificationSlideoverOpen } = useDashboard();

const notifications = useNotificationsStore();
const { notificationsCount, doNotDisturb } = storeToRefs(notifications);

const newNotification = ref(false);

const { start } = useTimeoutFn(() => (newNotification.value = false), 1000, {
    immediate: false,
});

const currentCount = ref(0);
watch(notificationsCount, () => {
    if (notificationsCount.value > currentCount.value) {
        newNotification.value = true;
        start();
    }

    currentCount.value = notificationsCount.value;
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
                src: '/images/logo.png',
                alt: 'FiveNet',
            }"
            :ui="{ base: collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)', leadingAvatar: 'rounded-none' }"
        >
            <span v-if="!collapsed" class="truncate font-semibold text-highlighted">FiveNet</span>
        </UButton>
    </UTooltip>

    <UTooltip v-if="!collapsed" :text="$t('components.partials.sidebar_notifications')" :kbds="['B']">
        <UChip
            :show="notificationsCount > 0"
            color="error"
            inset
            :text="notificationsCount <= 9 ? notificationsCount : '9+'"
            size="xl"
        >
            <UButton
                color="neutral"
                variant="ghost"
                square
                :icon="
                    doNotDisturb
                        ? 'i-mdi-bell-off-outline'
                        : notificationsCount === 0
                          ? 'i-mdi-notifications-none'
                          : 'i-mdi-notifications'
                "
                @click="isNotificationSlideoverOpen = true"
            />
        </UChip>
    </UTooltip>
</template>

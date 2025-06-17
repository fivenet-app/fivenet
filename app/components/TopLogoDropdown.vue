<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

const { isNotificationsSlideoverOpen } = useDashboard();

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
</script>

<template>
    <UDropdown v-slot="{ open }" class="w-full" :ui="{ width: 'w-full' }" :popper="{ strategy: 'absolute' }">
        <UButton class="w-full" :class="[open && 'bg-gray-50 dark:bg-gray-800']" color="gray" variant="ghost">
            <FiveNetLogo class="size-4" />

            <span class="truncate font-semibold text-gray-900 dark:text-white">FiveNet</span>
        </UButton>
    </UDropdown>

    <UTooltip :text="$t('components.partials.sidebar_notifications')" :shortcuts="['B']">
        <UChip
            :show="notificationsCount > 0"
            color="error"
            inset
            :text="notificationsCount <= 9 ? notificationsCount : '9+'"
            size="xl"
        >
            <UButton color="gray" variant="ghost" square @click="isNotificationsSlideoverOpen = true">
                <UIcon
                    class="size-5"
                    :name="
                        doNotDisturb
                            ? 'i-mdi-bell-off-outline'
                            : notificationsCount === 0
                              ? 'i-mdi-notifications-none'
                              : 'i-mdi-notifications'
                    "
                />
            </UButton>
        </UChip>
    </UTooltip>
</template>

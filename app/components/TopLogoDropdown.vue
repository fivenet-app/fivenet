<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import { useNotificatorStore } from '~/stores/notificator';

const { isNotificationsSlideoverOpen } = useDashboard();

const notificatorStore = useNotificatorStore();
const { notificationsCount, doNotDisturb } = storeToRefs(notificatorStore);

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
    <UDropdownMenu v-slot="{ open }" class="w-full" :ui="{ width: 'w-full' }" :popper="{ strategy: 'absolute' }">
        <UButton color="neutral" variant="ghost" :class="[open && 'bg-neutral-50 dark:bg-neutral-800']" class="w-full">
            <FiveNetLogo class="size-4" />

            <span class="truncate font-semibold text-neutral-900 dark:text-white">FiveNet</span>
        </UButton>
    </UDropdownMenu>

    <UTooltip :text="$t('components.partials.sidebar_notifications')" :shortcuts="['B']">
        <UChip
            :show="notificationsCount > 0"
            color="error"
            inset
            :text="notificationsCount <= 9 ? notificationsCount : '9+'"
            size="xl"
        >
            <UButton color="neutral" variant="ghost" square @click="isNotificationsSlideoverOpen = true">
                <UIcon
                    :name="
                        doNotDisturb
                            ? 'i-mdi-bell-off-outline'
                            : notificationsCount === 0
                              ? 'i-mdi-notifications-none'
                              : 'i-mdi-notifications'
                    "
                    class="size-5"
                />
            </UButton>
        </UChip>
    </UTooltip>
</template>

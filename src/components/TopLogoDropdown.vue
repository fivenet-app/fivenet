<script lang="ts" setup>
import { useNotificatorStore } from '~/store/notificator';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

const { isNotificationsSlideoverOpen } = useDashboard();

const notificatorStore = useNotificatorStore();
const { getNotificationsCount } = storeToRefs(notificatorStore);

const newNotification = ref(false);

const { start } = useTimeoutFn(() => (newNotification.value = false), 1000, {
    immediate: false,
});

const currentCount = ref(0);
watch(getNotificationsCount, () => {
    if (getNotificationsCount.value > currentCount.value) {
        newNotification.value = true;
        start();
    }

    currentCount.value = getNotificationsCount.value;
});
</script>

<template>
    <UDropdown v-slot="{ open }" class="w-full" :ui="{ width: 'w-full' }" :popper="{ strategy: 'absolute' }">
        <UButton color="gray" variant="ghost" :class="[open && 'bg-gray-50 dark:bg-gray-800']" class="w-full">
            <FiveNetLogo class="size-4" />

            <span class="truncate font-semibold text-gray-900 dark:text-white">FiveNet</span>
        </UButton>
    </UDropdown>

    <UTooltip text="Notifications" :shortcuts="['B']">
        <UChip
            :show="getNotificationsCount > 0"
            color="red"
            inset
            :text="getNotificationsCount <= 9 ? getNotificationsCount : '9+'"
            size="xl"
        >
            <UButton color="gray" variant="ghost" square @click="isNotificationsSlideoverOpen = true">
                <UIcon name="i-mdi-bell-outline" class="size-5" />
            </UButton>
        </UChip>
    </UTooltip>
</template>

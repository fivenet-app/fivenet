<script lang="ts" setup>
const notifications = useNotificationsStore();
const { doNotDisturb } = storeToRefs(notifications);
</script>

<template>
    <div class="inline-flex items-center gap-2">
        <USwitch
            v-model="doNotDisturb"
            size="lg"
            color="error"
            checked-icon="i-mdi-do-not-disturb-on"
            unchecked-icon="i-mdi-do-not-disturb-off-outline"
        />

        <span class="text-sm font-medium text-gray-700 dark:text-gray-200">
            {{ $t('components.notifications.do_not_disturb') }}
        </span>
    </div>
</template>

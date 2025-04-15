<script lang="ts" setup>
import { useNotificatorStore } from '~/stores/notificator';

const notifications = useNotificatorStore();
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

        <span class="text-sm font-medium text-neutral-700 dark:text-neutral-200">
            {{ $t('components.notifications.do_not_disturb') }}
        </span>
    </div>
</template>

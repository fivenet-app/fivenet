<script lang="ts" setup>
export type NotificationDeliverySetting = 'inboxEnabled' | 'toastEnabled' | 'soundEnabled';

const props = defineProps<{
    label?: string;
    settings: Record<NotificationDeliverySetting, boolean>;
    pending?: boolean;
    canReset?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update', setting: NotificationDeliverySetting, value: boolean): void;
    (e: 'reset'): void;
}>();

const settingLabels = {
    inboxEnabled: 'inbox',
    toastEnabled: 'toast',
    soundEnabled: 'sound',
} as const;

const settingKeys = Object.keys(settingLabels) as NotificationDeliverySetting[];
</script>

<template>
    <div :class="props.label ? 'grid gap-2 lg:grid-cols-[minmax(12rem,1fr)_auto_auto_auto_auto] lg:items-center' : undefined">
        <div v-if="props.label" class="flex items-center gap-2 pl-7">
            <span class="text-sm text-default">{{ props.label }}</span>
        </div>

        <div class="flex flex-wrap items-center gap-x-4 gap-y-2" :class="{ 'lg:col-span-4': props.label }">
            <UTooltip v-if="props.canReset" :text="$t('common.reset')">
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-restore"
                    class="!p-0"
                    size="xs"
                    :loading="props.pending"
                    :ui="{ leadingIcon: 'size-4' }"
                    @click.stop="emit('reset')"
                />
            </UTooltip>

            <UCheckbox
                v-for="setting in settingKeys"
                :key="setting"
                :model-value="props.settings[setting]"
                :label="$t(`components.auth.user_settings.notification_delivery.${settingLabels[setting]}`)"
                :disabled="props.pending"
                @click.stop
                @update:model-value="(value) => emit('update', setting, value === true)"
            />
        </div>
    </div>
</template>

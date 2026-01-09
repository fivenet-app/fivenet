<script lang="ts" setup>
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        value: string | number | (() => string | number);
        showText?: boolean;
    }>(),
    {
        showText: false,
    },
);

const notifications = useNotificationsStore();

function addToClipboard(): void {
    copyToClipboardWrapper(typeof props.value === 'function' ? props.value().toString() : props.value.toString());

    notifications.add({
        title: { key: 'notifications.clipboard.copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UTooltip :text="$t('components.clipboard.clipboard_button.add')">
        <UButton
            class="px-1 py-1"
            trailing-icon="i-mdi-clipboard-plus"
            variant="outline"
            color="neutral"
            size="xs"
            :label="props.showText ? $t('common.copy') : undefined"
            @click="addToClipboard()"
        />
    </UTooltip>
</template>

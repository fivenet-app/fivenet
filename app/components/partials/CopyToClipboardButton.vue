<script lang="ts" setup>
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    value: string | number;
}>();

const notifications = useNotificationsStore();

function addToClipboard(): void {
    copyToClipboardWrapper(props.value.toString());

    notifications.add({
        title: { key: 'notifications.clipboard.copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}
</script>

<template>
    <UTooltip :text="$t('components.clipboard.clipboard_button.add')">
        <UButton
            class="px-1 py-1"
            icon="i-mdi-clipboard-plus"
            variant="outline"
            color="neutral"
            size="xs"
            @click="addToClipboard()"
        />
    </UTooltip>
</template>

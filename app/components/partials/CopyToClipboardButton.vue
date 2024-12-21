<script lang="ts" setup>
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    value: string | number;
}>();

const notifications = useNotificatorStore();

function addToClipboard(): void {
    copyToClipboardWrapper(props.value.toString());

    notifications.add({
        title: { key: 'notifications.clipboard.copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.copied.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}
</script>

<template>
    <UTooltip :text="$t('components.clipboard.clipboard_button.add')">
        <UButton
            icon="i-mdi-clipboard-plus"
            variant="outline"
            color="black"
            size="xs"
            :ui="{ padding: { xs: '' } }"
            class="px-1 py-1"
            @click="addToClipboard()"
        />
    </UTooltip>
</template>

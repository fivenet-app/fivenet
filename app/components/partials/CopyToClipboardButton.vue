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
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}
</script>

<template>
    <UButton
        icon="i-mdi-clipboard-plus"
        variant="outline"
        color="black"
        size="xs"
        :ui="{ padding: { xs: '' } }"
        class="px-1 py-1"
        :title="$t('components.clipboard.clipboard_button.add')"
        @click="addToClipboard()"
    />
</template>

<script lang="ts" setup>
import type { ButtonSize, ButtonVariant } from '#ui/types';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        id: string | number | string;
        prefix?: string;
        title?: TranslateItem;
        content?: TranslateItem;
        action?: (id: string | number | string) => void;
        hideIcon?: boolean;
        variant?: ButtonVariant;
        padded?: boolean;
        size?: ButtonSize;
    }>(),
    {
        prefix: undefined,
        title: undefined,
        content: undefined,
        action: undefined,
        hideIcon: false,
        variant: 'solid',
        padded: true,
        size: 'sm',
    },
);

const notifications = useNotificatorStore();

function copyDocumentIDToClipboard(): void {
    copyToClipboardWrapper(props.prefix ? props.prefix + '-' + props.id : props.id.toString());

    if (props.title && props.content) {
        notifications.add({
            title: props.title,
            description: props.content,
            timeout: 3250,
            type: NotificationType.INFO,
        });
    }
}

function click(): void {
    if (props.action !== undefined) {
        props.action(props.id);
    } else {
        copyDocumentIDToClipboard();
    }
}
</script>

<template>
    <UButton
        :ui="{ round: 'rounded-md' }"
        :icon="!hideIcon ? 'i-mdi-fingerprint' : undefined"
        :variant="variant"
        :padded="padded"
        :size="size"
        class="break-keep"
        @click="click"
    >
        <template v-if="prefix">{{ prefix }}-</template>{{ id }}
    </UButton>
</template>

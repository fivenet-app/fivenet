<script lang="ts" setup>
import type { ButtonVariant } from '#ui/types';
import { type TranslateItem } from '~/composables/i18n';
import { useNotificatorStore } from '~/store/notificator';

const notifications = useNotificatorStore();

const props = withDefaults(
    defineProps<{
        id: string | number | string;
        prefix: string;
        title?: TranslateItem;
        content?: TranslateItem;
        action?: (id: string | number | string) => void;
        hideIcon?: boolean;
        variant?: ButtonVariant;
        padded?: boolean;
    }>(),
    {
        title: undefined,
        content: undefined,
        action: undefined,
        hideIcon: false,
        variant: 'solid',
        padded: true,
    },
);

function copyDocumentIDToClipboard(): void {
    copyToClipboardWrapper(props.prefix + '-' + props.id);

    if (props.title && props.content) {
        notifications.add({
            title: props.title,
            description: props.content,
            timeout: 3250,
            type: 'info',
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
        class="break-keep"
        @click="click"
    >
        {{ prefix }}-{{ id }}
    </UButton>
</template>

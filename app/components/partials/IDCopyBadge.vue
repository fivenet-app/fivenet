<script lang="ts" setup>
import type { ButtonSize, ButtonVariant } from '#ui/types';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        id: number | number | string;
        prefix?: string;
        title?: TranslateItem;
        content?: TranslateItem;
        action?: (id: number | number | string) => void;
        hideIcon?: boolean;
        variant?: ButtonVariant;
        padded?: boolean;
        size?: ButtonSize;
        disableTooltip?: boolean;
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
        disableTooltip: false,
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
    <UTooltip :text="disableTooltip ? undefined : $t('common.copy')">
        <UButton
            :ui="{ round: 'rounded-md', base: '' }"
            :icon="!hideIcon ? 'i-mdi-fingerprint' : undefined"
            :variant="variant"
            :padded="padded"
            :size="size"
            class="break-keep"
            @click="click"
        >
            <span class="hidden sm:block">
                <template v-if="prefix">{{ prefix }}-</template>{{ id }}
            </span>
        </UButton>
    </UTooltip>
</template>

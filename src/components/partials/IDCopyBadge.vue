<script lang="ts" setup>
import { type TranslateItem } from '~/composables/i18n';
import { useNotificatorStore } from '~/store/notificator';

const notifications = useNotificatorStore();

const props = defineProps<{
    id: string | number | string;
    prefix: string;
    title?: TranslateItem;
    content?: TranslateItem;
    action?: (id: string | number | string) => void;
    hideIcon?: boolean;
}>();

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
        :icon="hideIcon === undefined || !hideIcon ? 'i-mdi-fingerprint' : undefined"
        class="break-keep"
        @click="click"
    >
        {{ prefix }}-{{ id }}
    </UButton>
</template>

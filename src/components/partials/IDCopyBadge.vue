<script lang="ts" setup>
import { FingerprintIcon } from 'mdi-vue3';
import { type TranslateItem } from '~/composables/i18n';
import { useNotificatorStore } from '~/store/notificator';

const notifications = useNotificatorStore();

const props = defineProps<{
    id: string | number | string;
    prefix: string;
    title?: TranslateItem;
    content?: TranslateItem;
    action?: (id: string | number | string) => void;
}>();

function copyDocumentIDToClipboard(): void {
    copyToClipboardWrapper(props.prefix + '-' + props.id);

    if (props.title && props.content) {
        notifications.dispatchNotification({
            title: props.title,
            content: props.content,
            duration: 3250,
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
    <button
        type="button"
        class="inline-flex flex-initial flex-row items-center gap-1 rounded-full bg-base-500 px-2 py-1 text-base-100"
        @click.prevent="click"
    >
        <FingerprintIcon class="w-5 h-auto" aria-hidden="true" />
        <span class="break-keep text-sm font-medium text-base-100">{{ prefix }}-{{ id }}</span>
    </button>
</template>

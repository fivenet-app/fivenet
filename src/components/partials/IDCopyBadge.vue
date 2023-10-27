<script lang="ts" setup>
import { FingerprintIcon } from 'mdi-vue3';
import { type TranslateItem } from '~/composables/i18n';
import { useNotificatorStore } from '~/store/notificator';

const notifications = useNotificatorStore();

const props = defineProps<{
    id: bigint | number | string;
    prefix: string;
    title?: TranslateItem;
    content?: TranslateItem;
    action?: (id: bigint | number | string) => void;
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
    <div class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full text-base-100 bg-base-500" @click="click">
        <FingerprintIcon class="w-5 h-auto" aria-hidden="true" />
        <span class="text-sm font-medium text-base-100 break-keep">{{ prefix }}-{{ id }}</span>
    </div>
</template>

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
    <button
        type="button"
        class="inline-flex flex-initial flex-row items-center gap-1 rounded-full bg-gray-600 px-2 py-1 text-neutral"
        @click.prevent="click"
    >
        <FingerprintIcon v-if="hideIcon === undefined || !hideIcon" class="h-auto w-5" aria-hidden="true" />
        <span class="break-keep text-sm font-medium text-neutral">{{ prefix }}-{{ id }}</span>
    </button>
</template>

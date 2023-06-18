<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiFingerprint } from '@mdi/js';
import { useClipboard } from '@vueuse/core';
import { useNotificationsStore } from '~/store/notifications';
import { TranslateItem } from '~~/gen/ts/resources/common/i18n';

const clipboard = useClipboard();
const notifications = useNotificationsStore();

const props = defineProps<{
    id: bigint | string;
    prefix: string;
    title?: TranslateItem;
    content?: TranslateItem;
}>();

function copyDocumentIDToClipboard(): void {
    clipboard.copy(props.prefix + '-' + props.id);

    if (props.title && props.content) {
        notifications.dispatchNotification({
            title: props.title,
            content: props.content,
            duration: 3500,
            type: 'info',
        });
    }
}
</script>

<template>
    <div
        class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full text-base-100 bg-base-500"
        @click="copyDocumentIDToClipboard"
    >
        <SvgIcon class="w-5 h-auto" aria-hidden="true" type="mdi" :path="mdiFingerprint" />
        <span class="text-sm font-medium text-base-100">{{ prefix }}-{{ id }}</span>
    </div>
</template>

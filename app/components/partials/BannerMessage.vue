<script lang="ts" setup>
import { useNotificatorStore } from '~/stores/notificator';
import type { BannerMessage } from '~~/gen/ts/resources/rector/banner';

const props = defineProps<{
    message: BannerMessage;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const notificationStore = useNotificatorStore();
const { dismissedBannerMessageID } = storeToRefs(notificationStore);

const now = new Date();

const hide = computed(
    () =>
        dismissedBannerMessageID.value === props.message.id ||
        (props.message.expiresAt !== undefined && toDate(props.message.expiresAt).getTime() - now.getTime() < 0),
);

function onClose() {
    dismissedBannerMessageID.value = props.message.id;
    emit('close');
}
</script>

<template>
    <div v-if="!hide" class="fixed top-0 z-50 w-full">
        <UBanner color="primary" close @close="onClose">
            <template #title>
                <!-- eslint-disable-next-line vue/no-v-html -->
                <span v-html="props.message.title" />
            </template>
        </UBanner>
    </div>
</template>

<script lang="ts" setup>
import { useNotificatorStore } from '~/stores/notificator';
import type { BannerMessage } from '~~/gen/ts/resources/settings/banner';

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
        <div class="bg-primary-600 flex justify-between gap-1 p-2" :class="`bg-${message.color ?? 'primary'}-600`">
            <div class="flex flex-1 items-center justify-center gap-1">
                <UIcon class="size-6" :name="message.icon ?? 'i-mdi-information-outline'" />

                <!-- eslint-disable-next-line vue/no-v-html -->
                <p class="font-medium text-gray-900 dark:text-white" v-html="message.title"></p>
            </div>

            <UButton class="self-end" variant="link" icon="i-mdi-close" color="white" @click="onClose" />
        </div>
    </div>
</template>

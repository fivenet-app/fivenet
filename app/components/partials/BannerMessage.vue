<script lang="ts" setup>
import { useNotificatorStore } from '~/store/notificator';
import type { BannerMessage } from '~~/gen/ts/resources/rector/banner';

const props = defineProps<{
    message: BannerMessage;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const notificationStore = useNotificatorStore();
const { dismissedBannerMessageID } = storeToRefs(notificationStore);

const show = computed(
    () =>
        dismissedBannerMessageID.value !== props.message.id ||
        (props.message.createdAt !== undefined &&
            props.message.expiresAt !== undefined &&
            toDate(props.message.createdAt).getTime() - toDate(props.message.expiresAt).getTime() > 0),
);

function onClose() {
    dismissedBannerMessageID.value = props.message.id;
    emit('close');
}
</script>

<template>
    <div v-if="show" class="fixed top-0 z-50 w-full">
        <div class="bg-primary-600 flex justify-between gap-1 p-2" :class="`bg-${message.color ?? 'primary'}-600`">
            <div class="flex flex-1 items-center justify-center gap-1">
                <UIcon :name="message.icon ?? 'i-mdi-information-outline'" class="size-6" />

                <!-- eslint-disable-next-line vue/no-v-html -->
                <p class="font-medium text-gray-900 dark:text-white" v-html="message.title"></p>
            </div>

            <UButton variant="link" icon="i-mdi-close" color="white" class="self-end" @click="onClose" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';
import type { BannerMessage } from '~~/gen/ts/resources/settings/banner';

const props = defineProps<{
    message: BannerMessage;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const notificationStore = useNotificationsStore();
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

const color = computed(() => (props.message.color ?? 'primary') as ButtonProps['color']);

const appConfig = useAppConfig();
</script>

<template>
    <UBanner
        v-if="appConfig.system.bannerMessage && !hide"
        :icon="appConfig.system.bannerMessage.icon ?? 'i-mdi-information-outline'"
        :color="color"
        @close="onClose"
    >
        <template #title>
            <!-- eslint-disable-next-line vue/no-v-html -->
            <p class="font-medium text-highlighted" v-html="appConfig.system.bannerMessage.title"></p>
        </template>
    </UBanner>
</template>

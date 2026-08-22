<script lang="ts" setup>
import { useResizeObserver } from '@vueuse/core';
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
        (props.message.expiresAt && toDate(props.message.expiresAt).getTime() - now.getTime() < 0),
);

const bannerRef = useTemplateRef<{ $el: HTMLElement }>('bannerRef');
const bannerEl = computed(() => bannerRef.value?.$el ?? null);

const bannerMessageBottomOffsetVar = '--banner-message-bottom-offset';

function setBannerMessageBottomOffset(height: number): void {
    if (typeof document === 'undefined') return;

    document.documentElement.style.setProperty(bannerMessageBottomOffsetVar, `${height}px`);
}

function syncBannerMessageBottomOffset(): void {
    setBannerMessageBottomOffset(bannerEl.value?.clientHeight ?? 0);
}

function onClose() {
    if (props.message.expiresAt) dismissedBannerMessageID.value = props.message.id;

    emit('close');
}

const color = computed(() => (props.message.color ?? 'primary') as ButtonProps['color']);
const { system } = useAppConfig();

useResizeObserver(bannerEl, syncBannerMessageBottomOffset);

watch(bannerEl, syncBannerMessageBottomOffset, { immediate: true });

onBeforeUnmount(() => setBannerMessageBottomOffset(0));
</script>

<template>
    <UBanner
        v-if="system.bannerMessage && !hide"
        ref="bannerRef"
        :icon="system.bannerMessage.icon ?? 'i-mdi-information-outline'"
        :color="color"
        close
        :ui="{
            root: 'w-full pointer-events-auto',
            container: 'flex items-start justify-between gap-3 min-h-12 h-auto py-2 max-h-16 items-center',
            center: 'flex items-start gap-1.5 min-w-0 flex-1',
            left: 'lg:hidden',
            right: 'lg:flex-1 flex items-start justify-end',
            title: 'text-sm text-inverted font-medium whitespace-normal break-words',
        }"
        @close="onClose"
    >
        <template #title>
            <!-- eslint-disable-next-line vue/no-v-html -->
            <div v-html="system.bannerMessage.title"></div>
        </template>
    </UBanner>
</template>

<script lang="ts" setup>
import { NuxtImg } from '#components';
import type { AvatarProps } from '@nuxt/ui';
import { useSettingsStore } from '~/stores/settings';

const props = withDefaults(
    defineProps<{
        src?: string;
        alt?: string;
        text?: string;
        size?: AvatarProps['size'];
        noBlur?: boolean;
        enablePopup?: boolean;
        disableBlurToggle?: boolean;
        rounded?: boolean;
        imgClass?: string;
    }>(),
    {
        src: undefined,
        alt: '',
        text: '',
        size: 'lg',
        noBlur: undefined,
        enablePopup: false,
        disableBlurToggle: false,
        rounded: true,
        imgClass: '',
    },
);

const settings = useSettingsStore();
const { streamerMode } = storeToRefs(settings);

const visible = ref(props.noBlur || !streamerMode.value);

function toggleBlur(): void {
    if (props.disableBlurToggle || !streamerMode.value) {
        return;
    }

    if ((streamerMode.value && props.noBlur === undefined) || props.noBlur === false) {
        visible.value = !visible.value;
    }
}

const src = computed(() => {
    if (!props.src) {
        return props.src;
    }

    if (!props.src.startsWith('http') && !props.src.startsWith('/images') && !props.src.startsWith('/api/filestore')) {
        return `/api/filestore/${props.src}`;
    }

    return props.src;
});
</script>

<template>
    <UPopover>
        <UButton
            variant="link"
            :disabled="!src || !enablePopup"
            :ui="{ base: !src || !enablePopup ? 'disabled:cursor-default' : 'cursor-pointer' }"
        >
            <UAvatar
                :class="[visible ? '' : 'blur', imgClass]"
                :size="size"
                :src="src"
                :alt="alt"
                :text="text"
                :ui="{ root: rounded ? 'rounded-full' : 'rounded-sm' }"
                v-bind="$attrs"
            />
        </UButton>

        <template #content>
            <div class="p-4">
                <NuxtImg
                    class="h-96 max-w-full"
                    :class="[visible ? '' : 'blur', rounded && 'rounded-md']"
                    :src="src"
                    :alt="alt"
                    loading="lazy"
                    @click="toggleBlur()"
                />
            </div>
        </template>
    </UPopover>
</template>

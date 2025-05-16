<script lang="ts" setup>
import type { AvatarSize } from '#ui/types';
import { useSettingsStore } from '~/stores/settings';

const props = withDefaults(
    defineProps<{
        src?: string;
        alt?: string;
        text?: string;
        size?: AvatarSize;
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
</script>

<template>
    <template v-if="!src || !enablePopup">
        <UAvatar
            :as="NuxtImg"
            :size="size"
            :class="[visible ? '' : 'blur', imgClass]"
            :src="src"
            :alt="alt"
            :text="text"
            :ui="{ rounded: rounded ? 'rounded-full' : 'rounded' }"
            :img-class="imgClass"
            loading="lazy"
            @click="toggleBlur()"
        />
    </template>
    <UPopover v-else>
        <UButton variant="link" :padded="false">
            <UAvatar
                :as="NuxtImg"
                :size="size"
                :class="[visible ? '' : 'blur', imgClass]"
                :src="src"
                :alt="alt"
                :text="text"
                :ui="{ rounded: rounded ? 'rounded-full' : 'rounded' }"
                :img-class="imgClass"
                loading="lazy"
            />
        </UButton>

        <template #panel>
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

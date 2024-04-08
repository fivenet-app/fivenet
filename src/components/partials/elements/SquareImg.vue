<script lang="ts" setup>
import type { AvatarSize } from '#ui/types';
import { useSettingsStore } from '~/store/settings';

const props = withDefaults(
    defineProps<{
        url?: string;
        alt?: string;
        text?: string;
        size?: AvatarSize;
        noBlur?: boolean;
        enablePopup?: boolean;
    }>(),
    {
        url: undefined,
        alt: '',
        text: '',
        size: 'lg',
        noBlur: undefined,
        enablePopup: false,
    },
);

const settings = useSettingsStore();
const { streamerMode } = storeToRefs(settings);

const visible = ref(props.noBlur || !streamerMode.value);

function toggleBlur(): void {
    if ((streamerMode.value && props.noBlur === undefined) || props.noBlur === false) {
        visible.value = !visible.value;
    }
}
</script>

<template>
    <template v-if="!url || !enablePopup">
        <UAvatar :size="size" :class="[visible ? '' : 'blur']" :src="url" :alt="alt" :text="text" @click="toggleBlur()" />
    </template>
    <UPopover v-else>
        <UButton variant="link" :padded="false">
            <UAvatar :size="size" :class="[visible ? '' : 'blur']" :src="url" :alt="alt" :text="text" />
        </UButton>

        <template #panel>
            <div class="p-4">
                <img class="w-96 max-w-full rounded-md" :src="url" :alt="alt" />
            </div>
        </template>
    </UPopover>
</template>

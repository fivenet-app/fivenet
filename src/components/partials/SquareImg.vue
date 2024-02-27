<script lang="ts" setup>
import { imageSize, type imageSizes } from '~/components/partials/helpers';
import { useSettingsStore } from '~/store/settings';

const props = withDefaults(
    defineProps<{
        url?: string;
        text?: string;
        size?: imageSizes;
        rounded?: boolean;
        noBlur?: boolean;
    }>(),
    {
        url: undefined,
        text: '',
        size: 'lg',
        rounded: false,
        noBlur: false,
    },
);

const size = computed(() => imageSize(props.size));

const settings = useSettingsStore();
const { streamerMode } = storeToRefs(settings);

const visible = ref(streamerMode.value || props.noBlur);

function toggleBlur(): void {
    if (!props.noBlur) {
        visible.value = !visible.value;
    }
}
</script>

<template>
    <span
        class="flex items-center justify-center bg-gray-500 ring-2 ring-base-600"
        :class="[size, rounded ? 'rounded-full' : 'rounded-md']"
    >
        <span v-if="!url" class="font-medium leading-none text-white">
            <slot name="initials" />
        </span>
        <img
            v-else
            :class="[size, rounded ? 'rounded-full' : 'rounded-md', noBlur && visible ? '' : 'blur']"
            :src="url"
            :alt="text"
            @click="toggleBlur()"
        />
    </span>
</template>

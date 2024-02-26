<script lang="ts" setup>
import { imageSize, type imageSizes } from '~/components/partials/helpers';

const props = withDefaults(
    defineProps<{
        url?: string;
        text?: string;
        size?: imageSizes;
        rounded?: boolean;
    }>(),
    {
        url: undefined,
        text: '',
        size: 'lg',
        rounded: false,
    },
);

const size = computed(() => imageSize(props.size));
</script>

<template>
    <span
        class="flex items-center justify-center bg-gray-500 ring-2 ring-base-600"
        :class="[size, rounded ? 'rounded-full' : 'rounded-md']"
    >
        <span v-if="!url" class="font-medium leading-none text-white">
            <slot name="initials" />
        </span>
        <img v-else :class="[size, rounded ? 'rounded-full' : 'rounded-md']" :src="url" :alt="text" />
    </span>
</template>

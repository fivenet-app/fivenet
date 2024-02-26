<script lang="ts" setup>
import { imageSize, type imageSizes } from '~/components/partials/helpers';

const props = withDefaults(
    defineProps<{
        url?: string;
        text?: string;
        size?: imageSizes;
    }>(),
    {
        url: undefined,
        text: '',
        size: 'lg',
    },
);

const size = computed(() => imageSize(props.size));
</script>

<template>
    <span class="flex items-center justify-center rounded-sm bg-gray-500 ring-2 ring-base-600" :class="size">
        <span v-if="!url" class="font-medium leading-none text-white">
            {{ text }}
        </span>
        <img v-else class="rounded-sm" :class="size" :src="url" :alt="text" />
    </span>
</template>

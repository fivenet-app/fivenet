<script lang="ts" setup>
import { imageSize, type imageSizes } from '~/components/partials/helpers';

const props = withDefaults(
    defineProps<{
        url?: string;
        firstname: string;
        lastname: string;
        size?: imageSizes;
    }>(),
    {
        url: undefined,
        size: 'lg',
    },
);

const size = computed(() => imageSize(props.size));
</script>

<template>
    <span class="flex items-center justify-center rounded-full bg-gray-500 ring-2 ring-base-600" :class="size">
        <span v-if="!url" class="font-medium leading-none text-white">
            {{ getInitials(`${firstname} ${lastname}`) }}
        </span>
        <img v-else class="rounded-full" :class="size" :src="url" :alt="$t('common.avatar')" />
    </span>
</template>

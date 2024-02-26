<script lang="ts" setup>
import { type imageSizes } from '~/components/partials/helpers';
import SquareImg from '~/components/partials/SquareImg.vue';

withDefaults(
    defineProps<{
        url?: string;
        name: string;
        size?: imageSizes;
        rounded?: boolean;
    }>(),
    {
        url: undefined,
        size: 'lg',
        rounded: false,
    },
);
</script>

<template>
    <SquareImg :text="$t('common.avatar')" :size="size" :rounded="rounded">
        <template #initials>
            {{ getInitials(name) }}
        </template>
    </SquareImg>
</template>

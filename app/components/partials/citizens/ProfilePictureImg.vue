<script lang="ts" setup>
import type { AvatarSize } from '#ui/types';
import SquareImg from '~/components/partials/elements/SquareImg.vue';

withDefaults(
    defineProps<{
        src?: string;
        name: string;
        size?: AvatarSize;
        enablePopup?: boolean;
        noBlur?: boolean;
        alt?: string;
        disableBlurToggle?: boolean;
    }>(),
    {
        src: undefined,
        size: 'lg',
        enablePopup: false,
        noBlur: undefined,
        alt: undefined,
        disableBlurToggle: false,
    },
);
</script>

<template>
    <SquareImg
        :src="src"
        :size="size"
        :text="getInitials(name)"
        :alt="alt ?? $t('common.avatar')"
        :enable-popup="enablePopup"
        :no-blur="noBlur || src === undefined"
        :disable-blur-toggle="disableBlurToggle"
    />
</template>

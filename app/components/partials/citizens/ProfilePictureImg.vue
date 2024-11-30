<script lang="ts" setup>
import type { AvatarSize } from '#ui/types';
import GenericImg from '~/components/partials/elements/GenericImg.vue';

withDefaults(
    defineProps<{
        src?: string;
        name: string;
        size?: AvatarSize;
        enablePopup?: boolean;
        noBlur?: boolean;
        alt?: string;
        disableBlurToggle?: boolean;
        imgClass?: string;
    }>(),
    {
        src: undefined,
        size: 'lg',
        enablePopup: false,
        noBlur: undefined,
        alt: undefined,
        disableBlurToggle: false,
        imgClass: '',
    },
);
</script>

<template>
    <GenericImg
        :src="src"
        :size="size"
        :text="getInitials(name)"
        :alt="alt ?? $t('common.avatar')"
        :enable-popup="enablePopup"
        :no-blur="noBlur || src === undefined"
        :disable-blur-toggle="disableBlurToggle"
        :img-class="imgClass"
    />
</template>

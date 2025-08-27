<script lang="ts" setup>
import type { AvatarProps } from '@nuxt/ui';
import GenericImg from '~/components/partials/elements/GenericImg.vue';

withDefaults(
    defineProps<{
        src?: string;
        name: string;
        size?: AvatarProps['size'];
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
        :alt="alt ?? $t('common.profile_picture')"
        :enable-popup="enablePopup"
        :no-blur="noBlur || src === undefined"
        :disable-blur-toggle="disableBlurToggle"
        :img-class="imgClass"
    />
</template>

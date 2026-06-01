<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';

withDefaults(
    defineProps<{
        idx: number;
        moveUp: (index: number) => void;
        moveDown: (index: number) => void;
        orientation?: 'vertical' | 'horizontal';
        button?: ButtonProps;
    }>(),
    {
        orientation: 'vertical',
        button: undefined,
    },
);

const buttonDefault = computed<ButtonProps>(() => ({
    size: 'xs',
    variant: 'link',
    square: true,
}));

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UFieldGroup :orientation="orientation">
        <UButton
            :icon="orientation === 'vertical' ? 'i-mdi-arrow-up' : 'i-mdi-arrow-left'"
            v-bind="{
                ...buttonDefault,
                ...button,
            }"
            @click="moveUp(idx)"
        />

        <UButton
            :icon="orientation === 'vertical' ? 'i-mdi-arrow-down' : 'i-mdi-arrow-right'"
            v-bind="{
                ...buttonDefault,
                ...button,
            }"
            @click="moveDown(idx)"
        />
    </UFieldGroup>
</template>

<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';

const props = withDefaults(
    defineProps<{
        idx: number;
        moveUp: (index: number) => void;
        moveDown: (index: number) => void;
        orientation?: 'vertical' | 'horizontal';
        direction?: 'vertical' | 'horizontal';
        button?: ButtonProps;
        disableUp?: boolean;
        disableDown?: boolean;
    }>(),
    {
        orientation: 'vertical',
        direction: undefined,
        button: undefined,
        disableUp: false,
        disableDown: false,
    },
);

const buttonDefault = computed<ButtonProps>(() => ({
    size: 'xs',
    variant: 'link',
    square: true,
}));

const arrowDirection = computed(() => props.direction ?? props.orientation);

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UFieldGroup :orientation="props.orientation">
        <UButton
            :icon="arrowDirection === 'vertical' ? 'i-mdi-arrow-up' : 'i-mdi-arrow-left'"
            :disabled="disableUp"
            v-bind="{
                ...buttonDefault,
                ...button,
            }"
            @click="moveUp(idx)"
        />

        <UButton
            :icon="arrowDirection === 'vertical' ? 'i-mdi-arrow-down' : 'i-mdi-arrow-right'"
            :disabled="disableDown"
            v-bind="{
                ...buttonDefault,
                ...button,
            }"
            @click="moveDown(idx)"
        />
    </UFieldGroup>
</template>

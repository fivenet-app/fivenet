<script lang="ts" setup>
import { computed, useAttrs } from 'vue';
import type { ButtonProps, FieldGroupProps } from '@nuxt/ui';

const props = withDefaults(
    defineProps<
        ButtonProps & {
            handleClass?: string;
            orientation?: FieldGroupProps['orientation'];
        }
    >(),
    {
        handleClass: 'handle',
        orientation: 'horizontal',
    },
);

const attrs = useAttrs();

const forwardedProps = computed(() => {
    const { handleClass: _handleClass, orientation: _orientation, ...rest } = props;
    const definedProps = Object.fromEntries(Object.entries(rest).filter(([, value]) => value !== undefined));

    return {
        ...definedProps,
        ...attrs,
    } as ButtonProps;
});

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UTooltip :text="$t('common.draggable')">
        <UButton
            class="cursor-move select-none"
            :class="[props.handleClass]"
            type="button"
            variant="link"
            :icon="props.orientation === 'vertical' ? 'i-mdi-drag-vertical' : 'i-mdi-drag-horizontal'"
            tabindex="-1"
            v-bind="forwardedProps"
        />
    </UTooltip>
</template>

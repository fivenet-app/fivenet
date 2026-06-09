<script lang="ts" setup>
import type { ChipProps } from '@nuxt/ui';
import { backgroundColors, primaryColors, type PaletteColor } from '~/utils/color';

type ColorPickerColor = PaletteColor | 'primary';

const color = defineModel<string | undefined>({ default: 'primary' });

defineOptions({
    inheritAttrs: false,
});

const availableColorOptions = [
    { label: 'primary', chip: { color: 'primary' }, class: 'bg-primary-500 dark:bg-primary-400' },
    ...primaryColors,
    ...backgroundColors,
] as const;

const availableColorLabels = availableColorOptions.map((option) => option.label);

function isColorPickerColor(value: string | undefined): value is ColorPickerColor {
    return value !== undefined && availableColorLabels.includes(value as ColorPickerColor);
}

const selectColor = computed<ColorPickerColor | undefined>({
    get: () => (isColorPickerColor(color.value) ? color.value : 'primary'),
    set: (value) => {
        color.value = value;
    },
});
</script>

<template>
    <ClientOnly>
        <USelectMenu
            v-model="selectColor"
            :items="[...availableColorOptions]"
            :placeholder="$t('common.color')"
            value-key="label"
            v-bind="$attrs"
        >
            <template #leading="{ modelValue, ui }">
                <UChip
                    v-if="modelValue"
                    :class="ui.itemLeadingChip()"
                    :color="availableColorOptions.find((c) => c.label === modelValue)?.chip?.color || 'primary'"
                    inset
                    standalone
                    :size="ui.itemLeadingChipSize() as ChipProps['size']"
                />
            </template>
        </USelectMenu>
    </ClientOnly>
</template>

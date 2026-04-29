<script lang="ts" setup>
import type { ChipProps } from '@nuxt/ui';
import { backgroundColors, primaryColors } from '~/utils/color';

const color = defineModel<string | undefined>({ default: 'primary' });

defineOptions({
    inheritAttrs: false,
});

const availableColorOptions = [...primaryColors, ...backgroundColors];
</script>

<template>
    <ClientOnly>
        <USelectMenu
            v-model="color"
            :items="availableColorOptions"
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

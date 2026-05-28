<script lang="ts" setup>
import type { Component } from 'vue';
import {
    availableIcons,
    fallbackIconComponent as defaultFallbackIconComponent,
    fallbackIconName,
    type IconEntry,
} from './icons';
import { titleCase } from 'scule';

const props = withDefaults(
    defineProps<{
        color?: string;
        hexColor?: string;
        fallbackIcon?: Component | IconEntry;
        clear?: boolean;
    }>(),
    {
        color: undefined,
        hexColor: undefined,
        fallbackIcon: () => defaultFallbackIconComponent,
        clear: false,
    },
);

defineOptions({
    inheritAttrs: false,
});

const icon = defineModel<string | undefined>('modelValue');

function getItemName(item: unknown): string | undefined {
    if (typeof item !== 'object' || item === null || !('name' in item)) {
        return undefined;
    }

    const { name } = item as { name?: unknown };
    return typeof name === 'string' ? name : undefined;
}
</script>

<template>
    <ClientOnly>
        <USelectMenu
            v-model="icon"
            class="max-h-100"
            :items="availableIcons"
            :search-input="{ placeholder: $t('common.search_field') }"
            :filter-fields="['name', 'label']"
            label-key="label"
            value-key="name"
            virtualize
            :clear="props.clear"
            :ui="{ viewport: 'max-h-98' }"
            v-bind="$attrs"
        >
            <template v-if="icon" #default>
                <div class="inline-flex items-center gap-1">
                    <UIcon
                        class="size-5"
                        :name="convertComponentIconNameToDynamic(icon ?? fallbackIconName)"
                        :style="{ color: props.hexColor ?? `var(--color-${props.color ?? 'primary'}-400)` }"
                    />

                    <span class="truncate">{{ titleCase(icon ?? $t('common.unknown')) }}</span>
                </div>
            </template>

            <template #item="{ item }">
                <div class="inline-flex items-center gap-1">
                    <UIcon
                        class="size-5"
                        :name="convertComponentIconNameToDynamic(getItemName(item) ?? '')"
                        :style="{ color: props.hexColor ?? `var(--color-${props.color ?? 'primary'}-400)` }"
                    />

                    <span class="truncate">{{ titleCase(getItemName(item) ?? $t('common.unknown')) }}</span>
                </div>
            </template>
        </USelectMenu>
    </ClientOnly>
</template>

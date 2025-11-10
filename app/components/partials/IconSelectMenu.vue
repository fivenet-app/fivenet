<script lang="ts" setup>
import type { DefineComponent } from 'vue';
import { availableIcons, fallbackIcon as defaultIcon, type IconEntry } from './icons';

withDefaults(
    defineProps<{
        color?: string;
        hexColor?: string;
        fallbackIcon?: DefineComponent | IconEntry;
    }>(),
    {
        color: undefined,
        hexColor: undefined,
        fallbackIcon: () => defaultIcon,
    },
);

defineOptions({
    inheritAttrs: false,
});

const icon = defineModel<string | undefined>('modelValue');
</script>

<template>
    <ClientOnly>
        <USelectMenu
            v-model="icon"
            :items="availableIcons"
            :search-input="{ placeholder: $t('common.search_field') }"
            :filter-fields="['name']"
            value-key="name"
            virtualize
            v-bind="$attrs"
        >
            <template v-if="icon" #default>
                <div class="inline-flex items-center gap-1">
                    <component
                        :is="availableIcons.find((item) => item.name === icon)?.component ?? fallbackIcon.component"
                        class="size-5"
                        :style="{ color: hexColor ?? `var(--color-${color ?? 'primary'}-500)` }"
                    />

                    <span class="truncate">{{ camelCaseToTitleCase(icon ?? $t('common.unknown')) }}</span>
                </div>
            </template>

            <template #item-label="{ item }">
                <div class="inline-flex items-center gap-1">
                    <component
                        :is="item?.component"
                        class="size-5"
                        :style="{ color: hexColor ?? `var(--color-${color ?? 'primary'}-500)` }"
                    />

                    <span class="truncate">{{ camelCaseToTitleCase(item.name) }}</span>
                </div>
            </template>
        </USelectMenu>
    </ClientOnly>
</template>

<script lang="ts" setup>
import type { DefineComponent } from 'vue';
import { availableIcons, fallbackIcon as defaultIcon, type IconEntry } from './icons';

const props = withDefaults(
    defineProps<{
        modelValue: string | undefined;
        color?: string;
        fallbackIcon?: DefineComponent | IconEntry;
    }>(),
    {
        color: undefined,
        fallbackIcon: () => defaultIcon,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: string | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const icon = useVModel(props, 'modelValue', emit);

async function iconSearch(query: string): Promise<IconEntry[]> {
    // Remove spaces from query as icon names don't have spaces
    query = query.toLowerCase().replaceAll(' ', '').trim();
    let count = 0;
    return availableIcons.filter((icon) => {
        if (count < 35 && icon.name?.toLowerCase()?.startsWith(query)) {
            count++;
            return true;
        }
        return false;
    });
}
</script>

<template>
    <ClientOnly>
        <USelectMenu
            v-model="icon"
            :searchable="iconSearch"
            searchable-lazy
            :searchable-placeholder="$t('common.search_field')"
            value-attribute="name"
            v-bind="$attrs"
        >
            <template #label>
                <component
                    :is="availableIcons.find((item) => item.name === icon)?.component ?? fallbackIcon"
                    class="size-5"
                    :style="{ fill: color }"
                />
                <span class="truncate">{{ camelCaseToTitleCase(icon ?? $t('common.unknown')) }}</span>
            </template>
            <template #option="{ option }">
                <component :is="option?.component" class="size-5" :style="{ color: color }" />
                <span class="truncate">{{ camelCaseToTitleCase(option.name) }}</span>
            </template>
        </USelectMenu>
    </ClientOnly>
</template>

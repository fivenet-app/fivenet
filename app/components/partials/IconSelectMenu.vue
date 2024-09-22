<script lang="ts" setup>
import type { DefineComponent } from 'vue';
import { markerFallbackIcon, markerIcons } from '~/components/livemap/helpers';

const props = withDefaults(
    defineProps<{
        modelValue: string | undefined;
        color?: string;
        fallbackIcon?: DefineComponent;
    }>(),
    {
        color: undefined,
        fallbackIcon: markerFallbackIcon,
    },
);

const emits = defineEmits<{
    (e: 'update:modelValue', value: string | undefined): void;
}>();

const icon = useVModel(props, 'modelValue', emits);

async function markerIconSearch(query: string): Promise<DefineComponent[]> {
    // Remove spaces from query as icon names don't have spaces
    query = query.toLowerCase().replaceAll(' ', '').trim();
    let count = 0;
    return markerIcons.filter((icon) => {
        if (count < 35 && icon.name?.toLowerCase()?.startsWith(query)) {
            count++;
            return true;
        }
        return false;
    });
}
</script>

<template>
    <USelectMenu
        v-model="icon"
        :searchable="markerIconSearch"
        searchable-lazy
        :searchable-placeholder="$t('common.search_field')"
        value-attribute="name"
    >
        <template #label>
            <component
                :is="markerIcons.find((item) => item.name === icon) ?? fallbackIcon"
                class="size-5"
                :style="{ fill: color }"
            />
            <span class="truncate">{{ camelCaseToTitleCase(icon ?? $t('common.unknown')) }}</span>
        </template>
        <template #option="{ option }">
            <component :is="option" class="size-5" :style="{ color: color }" />
            <span class="truncate">{{ camelCaseToTitleCase(option.name) }}</span>
        </template>
    </USelectMenu>
</template>

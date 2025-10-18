<script lang="ts" setup>
import type { DefineComponent } from 'vue';
import { availableIcons, fallbackIcon as defaultIcon, type IconEntry } from './icons';

withDefaults(
    defineProps<{
        color?: string;
        fallbackIcon?: DefineComponent | IconEntry;
    }>(),
    {
        color: undefined,
        fallbackIcon: () => defaultIcon,
    },
);

defineOptions({
    inheritAttrs: false,
});

const icon = defineModel<string | undefined>('modelValue');

const searchTerm = ref('');
const searchTermDebounced = debouncedRef(searchTerm, 200);

async function iconSearch(query: string): Promise<IconEntry[]> {
    // Remove spaces from query as icon names don't have spaces
    query = query.toLowerCase().replaceAll(' ', '').trim();
    let count = 0;
    return availableIcons.filter((icon) => {
        if (count <= 40 && icon.name?.toLowerCase()?.includes(query)) {
            count++;
            return true;
        }
        return false;
    });
}

const foundIcons = computedAsync(() => iconSearch(searchTermDebounced.value));
</script>

<template>
    <ClientOnly>
        <USelectMenu
            v-model="icon"
            v-model:search-term="searchTerm"
            :items="foundIcons"
            :search-input="{ placeholder: $t('common.search_field') }"
            ignore-filter
            value-key="name"
            v-bind="$attrs"
        >
            <template v-if="icon" #default>
                <div class="inline-flex items-center gap-1">
                    <component
                        :is="availableIcons.find((item) => item.name === icon)?.component ?? fallbackIcon.component"
                        class="size-5"
                        :style="{ color: `var(--color-${color ?? 'primary'}-500)` }"
                    />

                    <span class="truncate">{{ camelCaseToTitleCase(icon ?? $t('common.unknown')) }}</span>
                </div>
            </template>

            <template #item-label="{ item }">
                <div class="inline-flex items-center gap-1">
                    <component
                        :is="item?.component"
                        class="size-5"
                        :style="{ color: `var(--color-${color ?? 'primary'}-500)` }"
                    />

                    <span class="truncate">{{ camelCaseToTitleCase(item.name) }}</span>
                </div>
            </template>
        </USelectMenu>
    </ClientOnly>
</template>

<script lang="ts" setup>
import { useSettingsStore, type LivemapLayer } from '~/stores/settings';
import { tileLayers } from '~/types/livemap';

const { attr, can } = useAuth();

const settingsStore = useSettingsStore();
const { livemapLayers, livemapLayerCategories, livemapTileLayer } = storeToRefs(settingsStore);

const groupedLayers = computed(() => {
    const reduced = livemapLayers.value.reduce(
        (grouped, l) => {
            if (
                l.perm === undefined ||
                (l.attr === undefined ? can(l.perm).value : attr(l.perm, l.attr.key, l.attr.val).value)
            ) {
                grouped[l.category] = (grouped[l.category] || []).concat(l);
            }
            return grouped;
        },
        {} as Record<string, LivemapLayer[]>,
    );
    // Sort the layers within each category
    Object.keys(reduced).forEach((key) => reduced[key]?.sort((a, b) => a.label.localeCompare(b.label)));

    // Reduce object to array and sort the categories by their order
    return Object.keys(reduced)
        .map((key) => ({
            category: livemapLayerCategories.value.find((c) => c.key === key),
            layers: reduced[key]!,
        }))
        .sort((a, b) => {
            if (a.category?.order !== undefined && b.category?.order !== undefined) {
                return a.category.order - b.category.order;
            }
            return 0;
        });
});
</script>

<template>
    <LControl position="topright">
        <UTooltip :text="$t('common.layer', 2)" :popper="{ placement: 'left' }">
            <UPopover :popper="{ arrow: true }" :ui="{ arrow: { base: 'mt-2' } }">
                <UButton
                    class="border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                    size="xl"
                    icon="i-mdi-layers-triple"
                    :ui="{ icon: { base: '!size-8' } }"
                />

                <template #panel>
                    <div class="w-full max-w-sm divide-y divide-gray-100 py-1 dark:divide-gray-800">
                        <div class="px-1">
                            <p class="truncate text-sm font-bold text-gray-900 dark:text-white">
                                {{ $t('common.layer', 2) }}
                            </p>

                            <URadioGroup
                                v-model="livemapTileLayer"
                                class="overflow-y-hidden"
                                :options="tileLayers"
                                value-attribute="key"
                                :ui-radio="{ inner: 'ms-1' }"
                                :ui="{ fieldset: 'grid auto-cols-auto grid-flow-col gap-1' }"
                            >
                                <template #label="{ option }">
                                    <span class="truncate">{{ $t(option.label) }}</span>
                                </template>
                            </URadioGroup>
                        </div>

                        <p v-if="Object.keys(groupedLayers).length === 0" class="truncate">
                            {{ $t('common.layers', 0) }}
                        </p>
                        <div v-else class="grid auto-cols-auto grid-flow-col divide-x divide-gray-100 dark:divide-gray-800">
                            <div
                                v-for="(category, key) in groupedLayers"
                                :key="key"
                                class="grid min-w-0 grid-flow-row auto-rows-min gap-1 px-1"
                            >
                                <p class="truncate text-sm font-bold text-gray-900 dark:text-white">
                                    {{ category.category?.label ?? $t('common.na') }}
                                </p>

                                <div
                                    v-for="layer in category.layers"
                                    :key="layer.key"
                                    class="inline-flex gap-1 overflow-y-hidden"
                                >
                                    <UToggle v-model="layer.visible" :disabled="!!layer.disabled" />
                                    <span class="truncate hover:line-clamp-2">{{ layer.label }}</span>
                                </div>
                            </div>
                        </div>

                        <slot />
                    </div>
                </template>
            </UPopover>
        </UTooltip>
    </LControl>
</template>

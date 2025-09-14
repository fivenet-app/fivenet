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
    const out = Object.keys(reduced)
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

    out.forEach((category) => {
        // Sort the layers within each category by their order
        category.layers.sort((a, b) => {
            if (a.order !== undefined && b.order !== undefined) {
                return a.order - b.order;
            }
            return 0;
        });
    });

    return out;
});
</script>

<template>
    <LControl position="topright">
        <UTooltip :text="$t('common.layer', 2)">
            <UPopover :ui="{ content: 'w-full' }">
                <UButton
                    class="border border-black/20 bg-clip-padding p-1.5"
                    size="xl"
                    icon="i-mdi-layers-triple"
                    :ui="{ leadingIcon: 'size-8!' }"
                />

                <template #content>
                    <div class="w-full max-w-xl divide-y divide-default py-1">
                        <div class="px-1 pb-0.5">
                            <p class="truncate text-base font-bold text-highlighted">
                                {{ $t('common.layer', 2) }}
                            </p>

                            <URadioGroup
                                v-model="livemapTileLayer"
                                class="overflow-y-hidden"
                                :items="tileLayers"
                                value-key="key"
                                :ui-radio="{ inner: 'ms-1' }"
                                :ui="{ fieldset: 'grid auto-cols-auto grid-flow-col gap-1' }"
                            >
                                <template #label="{ item }">
                                    {{ $t(item.label) }}
                                </template>
                            </URadioGroup>
                        </div>

                        <p v-if="Object.keys(groupedLayers).length === 0" class="truncate">
                            {{ $t('common.layers', 0) }}
                        </p>
                        <div
                            v-else
                            class="flex auto-cols-auto grid-flow-col flex-col divide-x divide-y divide-default md:grid md:divide-y-0"
                        >
                            <div
                                v-for="(category, key) in groupedLayers"
                                :key="key"
                                class="grid min-w-0 grid-flow-row auto-rows-min gap-1 overflow-y-hidden px-1 pb-1 md:pb-0"
                            >
                                <p class="truncate text-base font-bold text-highlighted">
                                    {{ category.category?.label ?? $t('common.na') }}
                                </p>

                                <USwitch
                                    v-for="layer in category.layers"
                                    :key="layer.key"
                                    v-model="layer.visible"
                                    :label="layer.label"
                                    :disabled="!!layer.disabled"
                                    :ui="{ label: 'truncate text-sm hover:line-clamp-2' }"
                                />
                            </div>
                        </div>

                        <slot />
                    </div>
                </template>
            </UPopover>
        </UTooltip>
    </LControl>
</template>

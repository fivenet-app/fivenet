<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { useSettingsStore, type LivemapLayer } from '~/stores/settings';

const { attr, can } = useAuth();

const settingsStore = useSettingsStore();
const { livemapLayers, livemapLayerCategories } = storeToRefs(settingsStore);

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

    Object.keys(reduced).forEach((key) => reduced[key]?.sort((a, b) => a.label.localeCompare(b.label)));

    return reduced;
});
</script>

<template>
    <LControl position="topright">
        <UTooltip :text="$t('common.layer', 2)" :popper="{ placement: 'left' }">
            <UPopover :popper="{ arrow: true }" :ui="{ arrow: { base: 'mt-2' } }">
                <UButton
                    size="xl"
                    icon="i-mdi-layers-triple"
                    class="border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                    :ui="{ icon: { base: '!size-8' } }"
                />

                <template #panel>
                    <div class="w-full max-w-sm py-1">
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
                                    {{ livemapLayerCategories.find((c) => c.key === key)?.label ?? $t('common.na') }}
                                </p>

                                <div v-for="layer in category" :key="layer.key" class="inline-flex gap-1 overflow-y-hidden">
                                    <UToggle v-model="layer.visible" />
                                    <span class="truncate hover:line-clamp-2">{{ layer.label }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </template>
            </UPopover>
        </UTooltip>
    </LControl>
</template>

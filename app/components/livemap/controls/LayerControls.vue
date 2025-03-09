<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { useSettingsStore, type LivemapLayer } from '~/store/settings';

const { attr, can } = useAuth();

const settingsStore = useSettingsStore();
const { livemapLayers } = storeToRefs(settingsStore);

const groupedLayers = useArrayReduce(
    livemapLayers,
    (grouped, l) => {
        if (l.perm === undefined || (l.attr === undefined ? can(l.perm).value : attr(l.perm, l.attr.key, l.attr.val).value)) {
            grouped[l.category] = (grouped[l.category] || []).concat(l);
        }
        return grouped;
    },
    {} as Record<string, LivemapLayer[]>,
);
</script>

<template>
    <LControl position="topright">
        <UPopover mode="hover" :popper="{ arrow: true }" :close-delay="150" :ui="{ arrow: { base: 'mt-2' } }">
            <UButton
                size="xl"
                icon="i-mdi-layers-triple"
                class="border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                :ui="{ icon: { base: '!size-8' } }"
            />

            <template #panel>
                <div class="flex flex-col gap-1 divide-y divide-gray-100 p-2 dark:divide-gray-800">
                    <p v-if="Object.keys(groupedLayers).length === 0">{{ $t('common.layers', 0) }}</p>

                    <template v-else>
                        <p class="text-sm font-medium text-gray-900 dark:text-white">{{ $t('common.layer', 2) }}</p>
                        <div v-for="(layers, lidx) in groupedLayers" :key="lidx" class="flex flex-col gap-1">
                            <div
                                v-for="(layer, idx) in layers"
                                :key="layer.key"
                                class="inline-flex items-center gap-1"
                                :class="idx === 0 && 'mt-1'"
                            >
                                <UToggle v-model="layer.visible" />
                                <span>{{ layer.label }}</span>
                            </div>
                        </div>
                    </template>
                </div>
            </template>
        </UPopover>
    </LControl>
</template>

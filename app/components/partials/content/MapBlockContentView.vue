<script lang="ts" setup>
import MapFullscreenModal from '~/components/livemap/MapFullscreenModal.vue';
import MapPositionPicker from '~/components/livemap/MapPositionPicker.vue';
import { mapTileLayers } from '~/composables/livemap/useMapProjection';
import { tileLayers } from '~/types/livemap';

const props = defineProps<{
    x: number;
    y: number;
    zoom: number;
    postal?: string;
    layer?: string;
    showGotoCoords?: boolean;
}>();

const open = ref(false);
const fullscreenOpen = ref(false);

const displayCoords = computed(() => `${props.x.toFixed(2)}, ${props.y.toFixed(2)}`);
const activeLayer = computed(() => mapTileLayers.find((layer) => layer.key === props.layer) ?? tileLayers[0]!);
const showGotoCoords = computed(() => props.showGotoCoords !== false);

function openFullscreen(): void {
    open.value = false;
    fullscreenOpen.value = true;
}
</script>

<template>
    <UPopover v-model:open="open">
        <UFieldGroup>
            <UButton
                class="overflow-hidden text-left text-sm"
                icon="i-mdi-map"
                color="neutral"
                size="sm"
                variant="outline"
                :ui="{ leadingIcon: 'size-7' }"
            >
                <div class="flex items-start gap-3 p-1">
                    <div class="min-w-0 flex-1">
                        <div class="font-medium">{{ postal || $t('common.map') }}</div>
                        <div class="mt-0.5 truncate text-xs text-neutral-500 dark:text-neutral-400">
                            {{ displayCoords }} · {{ $t(activeLayer.label) }} · {{ $t('common.zoom') }} {{ zoom }}
                        </div>
                    </div>
                </div>
            </UButton>

            <UTooltip v-if="showGotoCoords" :text="$t('common.goto')">
                <UButton
                    class="!my-0"
                    size="sm"
                    icon="i-mdi-map-marker"
                    color="neutral"
                    variant="subtle"
                    :to="`/livemap?loc=${zoom}/${y}/${x}`"
                    @click.prevent="
                        () => {
                            navigateTo(`/livemap?loc=${zoom}/${y}/${x}`);
                        }
                    "
                />
            </UTooltip>
        </UFieldGroup>

        <template #content>
            <div class="flex w-[38rem] max-w-[calc(100vw-2rem)] flex-col gap-3 p-4">
                <div class="flex items-start justify-between gap-4">
                    <div>
                        <div class="font-medium">{{ postal || $t('common.map') }}</div>
                        <div class="text-xs text-neutral-500 dark:text-neutral-400">
                            {{ displayCoords }} · {{ $t(activeLayer.label) }} · {{ $t('common.zoom') }} {{ zoom }}
                        </div>
                    </div>

                    <UFieldGroup>
                        <UTooltip :text="$t('common.fullscreen_enter')">
                            <UButton
                                color="neutral"
                                variant="ghost"
                                icon="i-mdi-fullscreen"
                                :aria-label="$t('common.fullscreen_enter')"
                                @click="openFullscreen"
                            />
                        </UTooltip>

                        <UButton
                            color="neutral"
                            variant="ghost"
                            trailing-icon="i-mdi-close"
                            :label="$t('common.close')"
                            @click="open = false"
                        />
                    </UFieldGroup>
                </div>

                <ClientOnly>
                    <MapPositionPicker v-if="open" :x="x" :y="y" :zoom="zoom" :layer="layer" disabled :show-controls="false" />
                </ClientOnly>
            </div>
        </template>
    </UPopover>

    <MapFullscreenModal
        v-model:open="fullscreenOpen"
        :title="postal || $t('common.map')"
        :summary="`${displayCoords} · ${$t(activeLayer.label)} · ${$t('common.zoom')} ${zoom}`"
        :x="x"
        :y="y"
        :zoom="zoom"
        :layer="layer"
        disabled
    />
</template>

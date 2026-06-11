<script lang="ts" setup>
import { createMapBlockAttrs, defaultMapBlockLayerKey, type MapBlockAttrs } from '~/composables/tiptap/extensions/MapBlock';
import MapFullscreenModal from '~/components/livemap/MapFullscreenModal.vue';
import MapPositionPicker from '~/components/livemap/MapPositionPicker.vue';
import PostalSearchSelect from '~/components/livemap/controls/PostalSearchSelect.vue';
import type { Postal, TileLayerKeys } from '~/types/livemap';

const props = defineProps<{
    modelValue: MapBlockAttrs;
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: MapBlockAttrs): void;
}>();

const draft = ref<MapBlockAttrs>(createMapBlockAttrs(props.modelValue));
const selectedPostal = ref<Postal | undefined>();
const fullscreenOpen = ref(false);

watch(
    () => props.modelValue,
    (value) => {
        draft.value = createMapBlockAttrs(value);
    },
    { immediate: true, deep: true },
);

function applyDraft(partial: Partial<MapBlockAttrs>): void {
    draft.value = {
        ...draft.value,
        ...partial,
    };

    emit('update:modelValue', draft.value);
}

watch(selectedPostal, (postal) => {
    if (!postal) return;
    if (postal.code === draft.value.postal && postal.x === draft.value.x && postal.y === draft.value.y) return;

    applyDraft({
        x: postal.x,
        y: postal.y,
        postal: postal.code,
    });
});

function applyPosition(x: number, y: number, zoom: number, layer: TileLayerKeys = defaultMapBlockLayerKey): void {
    selectedPostal.value = undefined;
    applyDraft({
        x,
        y,
        zoom,
        postal: '',
        layer,
    });
}

function openFullscreen(): void {
    fullscreenOpen.value = true;
}
</script>

<template>
    <div class="flex flex-col gap-3">
        <UFormField :label="$t('common.postal')" name="postal">
            <PostalSearchSelect v-model="selectedPostal" :selected-code="draft.postal" :disabled="disabled" />
        </UFormField>

        <div class="flex flex-col gap-1">
            <div class="flex items-center justify-between gap-2 text-xs text-neutral-500 dark:text-neutral-400">
                <span>{{ $t('common.select') }}</span>
                <div class="inline-flex items-center gap-2">
                    <span>{{ draft.x.toFixed(2) }}, {{ draft.y.toFixed(2) }} · {{ $t('common.zoom') }} {{ draft.zoom }}</span>
                    <UTooltip :text="$t('common.fullscreen_enter')">
                        <UButton
                            color="neutral"
                            variant="ghost"
                            size="xs"
                            icon="i-mdi-fullscreen"
                            :aria-label="$t('common.fullscreen_enter')"
                            @click="openFullscreen"
                        />
                    </UTooltip>
                </div>
            </div>

            <MapPositionPicker
                :x="draft.x"
                :y="draft.y"
                :zoom="draft.zoom"
                :layer="draft.layer ?? defaultMapBlockLayerKey"
                :disabled="disabled"
                @update:x="(value) => applyPosition(value, draft.y, draft.zoom, draft.layer)"
                @update:y="(value) => applyPosition(draft.x, value, draft.zoom, draft.layer)"
                @update:zoom="(value) => applyPosition(draft.x, draft.y, value, draft.layer)"
                @update:layer="(value) => applyDraft({ layer: value })"
            />
        </div>
    </div>

    <MapFullscreenModal
        v-model:open="fullscreenOpen"
        :title="$t('common.map')"
        :summary="`${draft.x.toFixed(2)}, ${draft.y.toFixed(2)} · ${$t('common.zoom')} ${draft.zoom}`"
        :x="draft.x"
        :y="draft.y"
        :zoom="draft.zoom"
        :layer="draft.layer ?? defaultMapBlockLayerKey"
        :disabled="disabled"
        @update:x="(value) => applyPosition(value, draft.y, draft.zoom, draft.layer)"
        @update:y="(value) => applyPosition(draft.x, value, draft.zoom, draft.layer)"
        @update:zoom="(value) => applyPosition(draft.x, draft.y, value, draft.layer)"
        @update:layer="(value) => applyDraft({ layer: value })"
    />
</template>

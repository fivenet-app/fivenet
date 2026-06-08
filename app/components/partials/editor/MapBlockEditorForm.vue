<script lang="ts" setup>
import type { MapBlockAttrs } from '~/composables/tiptap/extensions/MapBlock';
import { deepToRaw } from '~/utils/deepToRaw';
import MapPositionPicker from '~/components/livemap/MapPositionPicker.vue';
import PostalSearchSelect from '~/components/livemap/controls/PostalSearchSelect.vue';
import { tileLayers, type Postal } from '~/types/livemap';

const props = defineProps<{
    modelValue: MapBlockAttrs;
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: MapBlockAttrs): void;
}>();

const draft = ref<MapBlockAttrs>(structuredClone(deepToRaw(props.modelValue)));
const selectedPostal = ref<Postal | undefined>();

watch(
    () => props.modelValue,
    (value) => {
        draft.value = structuredClone(deepToRaw(value));
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

function applyPosition(x: number, y: number, zoom: number, layer: string): void {
    selectedPostal.value = undefined;
    applyDraft({
        x,
        y,
        zoom,
        postal: '',
        layer,
    });
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
                <span>{{ draft.x.toFixed(2) }}, {{ draft.y.toFixed(2) }} · {{ $t('common.zoom') }} {{ draft.zoom }}</span>
            </div>

            <MapPositionPicker
                :x="draft.x"
                :y="draft.y"
                :zoom="draft.zoom"
                :layer="draft.layer ?? ''"
                :disabled="disabled"
                @update:x="(value) => applyPosition(value, draft.y, draft.zoom, draft.layer || tileLayers[0]!.key)"
                @update:y="(value) => applyPosition(draft.x, value, draft.zoom, draft.layer || tileLayers[0]!.key)"
                @update:zoom="(value) => applyPosition(draft.x, draft.y, value, draft.layer || tileLayers[0]!.key)"
                @update:layer="(value) => applyDraft({ layer: value })"
            />
        </div>
    </div>
</template>

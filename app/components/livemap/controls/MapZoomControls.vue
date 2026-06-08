<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        modelValue: number;
        minZoom?: number;
        maxZoom?: number;
        disabled?: boolean;
        position?: 'topleft' | 'topright' | 'bottomleft' | 'bottomright';
    }>(),
    {
        minZoom: 1,
        maxZoom: 7,
        disabled: false,
        position: 'topleft',
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
}>();

const zoom = computed({
    get: () => props.modelValue,
    set: (value: number) => emit('update:modelValue', value),
});

function clampZoom(value: number): number {
    return Math.min(props.maxZoom, Math.max(props.minZoom, value));
}

function updateZoom(nextZoom: number): void {
    if (props.disabled) return;
    zoom.value = clampZoom(nextZoom);
}
</script>

<template>
    <LControl :position="position">
        <UFieldGroup class="inline-flex w-full flex-col" orientation="vertical">
            <UTooltip :text="$t('common.zoom_in')">
                <UButton
                    class="inset-0 border border-black/20 bg-clip-padding p-1.5"
                    icon="i-mdi-plus-thick"
                    block
                    :disabled="disabled || zoom >= maxZoom"
                    @click="updateZoom(zoom + 1)"
                />
            </UTooltip>

            <UTooltip :text="$t('common.zoom')">
                <UButton class="inset-0 border border-black/20 bg-clip-padding select-none" :label="zoom.toString()" block />
            </UTooltip>

            <UTooltip :text="$t('common.zoom_out')">
                <UButton
                    class="inset-0 border border-black/20 bg-clip-padding p-1.5"
                    icon="i-mdi-minus-thick"
                    block
                    :disabled="disabled || zoom <= minZoom"
                    @click="updateZoom(zoom - 1)"
                />
            </UTooltip>
        </UFieldGroup>
    </LControl>
</template>

<script lang="ts" setup>
import type { PopupOptions } from 'leaflet';

withDefaults(
    defineProps<{
        options?: PopupOptions;
    }>(),
    {
        options: () => ({ closeButton: false }),
    },
);

const popupRef = useTemplateRef<{ leafletObject?: { close?: () => void } }>('popupRef');
const popupOpen = ref<boolean>(false);

function closePopup(): void {
    popupRef.value?.leafletObject?.close?.();
}
</script>

<template>
    <LPopup ref="popupRef" :options="options" @add="popupOpen = true" @remove="popupOpen = false">
        <slot v-bind="{ popupOpen, closePopup }" />
    </LPopup>
</template>

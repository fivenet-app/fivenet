<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
    api?: any | null;
}>();

const zoom = defineModel<number>('zoom');

const zoomDisplay = computed(() => Math.round((zoom.value ?? 1) * 100));

const modes = [
    { value: 'select', label: 'Select' },
    { value: 'pan', label: 'Pan' },
    { value: 'path', label: 'Path' },
    { value: 'rect', label: 'Rect' },
    { value: 'ellipse', label: 'Ellipse' },
    { value: 'line', label: 'Line' },
    { value: 'text', label: 'Text' },
];

async function exportSvg(): Promise<void> {
    const svg = await props.api?.getSvgString?.();
    const blob = new Blob([svg], { type: 'image/svg+xml' });
    const url = URL.createObjectURL(blob);
    const a = doc.createElement('a');
    a.href = url;
    a.download = 'drawing.svg';
    a.click();
    URL.revokeObjectURL(url);
}

const doc = document;
</script>

<template>
    <div class="flex items-center gap-2 rounded-t-xl border-b p-2">
        <USelect
            :items="modes"
            option-attribute="label"
            value-attribute="value"
            size="sm"
            placeholder="Mode"
            @update:model-value="(v) => api?.setMode?.(v as string)"
        />

        <UInput type="color" size="sm" @change="(e) => api?.setColor?.('fill', (e.target as HTMLInputElement).value)" />
        <UInput type="color" size="sm" @change="(e) => api?.setColor?.('stroke', (e.target as HTMLInputElement).value)" />
        <UInput
            type="number"
            size="sm"
            placeholder="Stroke"
            class="w-20"
            @change="(e) => api?.setStrokeWidth?.(Number((e.target as HTMLInputElement).value || 1))"
        />

        <UButton size="sm" @click="api?.undoMgr?.undo?.()">Undo</UButton>
        <UButton size="sm" @click="api?.undoMgr?.redo?.()">Redo</UButton>

        <UButton size="sm" @click="zoom = (api?.zoom ?? 1) * 0.9">-</UButton>
        <span class="w-10 text-center text-xs">{{ zoomDisplay }}%</span>
        <UButton size="sm" variant="ghost" @click="api?.updateCanvas(true)">FIT</UButton>
        <UButton size="sm" @click="zoom = (api?.zoom ?? 1) * 1.1">+</UButton>

        <div class="flex-1" />

        <UButton size="sm" variant="ghost" @click="exportSvg">Export</UButton>

        <UButton size="sm" variant="outline" @click="() => api?.clear?.()">Clear</UButton>
    </div>
</template>

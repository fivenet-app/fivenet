<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import { useToast } from '#imports';
import { computed, h, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import FrameNode from './FrameNode.vue';
import Ruler from './Ruler.vue';
import type {
    BaseFrame,
    CheckboxFrame,
    FieldFrame,
    Frame,
    GridFrame,
    ImageFrame,
    Kind,
    LineFrame,
    RepeatFrame,
    SectionFrame,
    TextFrame,
} from './types';

/* Units */
const DPI = 96;
const mmToPx = (mm: number) => (mm / 25.4) * DPI;
const pxToMm = (px: number) => (px / DPI) * 25.4;

/* Page state */
const page = reactive({
    size: 'A4 portrait',
    widthMm: 210,
    heightMm: 297,
    backgroundUrl: '',
    bgColor: '#ffffff',
    bgOpacity: 0.45,
    bgLayer: 'behind' as 'behind' | 'overlay',
    bgBlend: 'normal' as any,
    bgLocked: false,
});

const pageSizeOptions = [
    { label: 'A4 portrait', value: 'A4 portrait' },
    { label: 'A4 landscape', value: 'A4 landscape' },
    { label: 'A5 portrait', value: 'A5 portrait' },
    { label: 'A5 landscape', value: 'A5 landscape' },
];
watch(
    () => page.size,
    (s) => {
        if (s === 'A4 portrait') {
            page.widthMm = 210;
            page.heightMm = 297;
        }
        if (s === 'A4 landscape') {
            page.widthMm = 297;
            page.heightMm = 210;
        }
        if (s === 'A5 portrait') {
            page.widthMm = 148;
            page.heightMm = 210;
        }
        if (s === 'A5 landscape') {
            page.widthMm = 210;
            page.heightMm = 148;
        }
    },
);
const bgLayerOptions = [
    { label: 'Behind content', value: 'behind' },
    { label: 'Overlay on top', value: 'overlay' },
];
const bgBlendOptions = ['normal', 'multiply', 'screen', 'overlay', 'darken', 'lighten'].map((v) => ({ label: v, value: v }));
const bgInput = ref<HTMLInputElement | null>(null);
function pickBg() {
    bgInput.value?.click();
}
function onBgChange(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;
    page.backgroundUrl = URL.createObjectURL(file);
}

/* Canvas */
const zoom = ref(0.9);
const gridStepMm = ref(2);
const gridStepOptions = [
    { label: '1 mm', value: 1 },
    { label: '2 mm', value: 2 },
    { label: '5 mm', value: 5 },
];
const imageFitOptions = [
    { label: 'Contain', value: 'contain' },
    { label: 'Cover', value: 'cover' },
    { label: 'Stretch', value: 'stretch' },
];
const kindOptions = [
    { label: 'Text', value: 'text' },
    { label: 'Field', value: 'field' },
    { label: 'Image', value: 'image' },
    { label: 'Checkbox', value: 'checkbox' },
    { label: 'Grid', value: 'grid' },
    { label: 'Line', value: 'line' },
    { label: 'Rotated Text', value: 'rotatedText' },
    { label: 'Section', value: 'section' },
];
const insertItems = [
    { label: 'Text', kind: 'text' },
    { label: 'Field', kind: 'field' },
    { label: 'Image', kind: 'image' },
    { label: 'Checkbox', kind: 'checkbox' },
    { label: 'Grid', kind: 'grid' },
    { label: 'Line', kind: 'line' },
    { label: 'Section', kind: 'section' },
] as const;

// Default for new frames: strokeEnabled true

const snap = ref(true);
const pagePx = computed(() => ({ width: mmToPx(page.widthMm), height: mmToPx(page.heightMm) }));
const pageStyle = computed(() => ({
    width: pagePx.value.width + 'px',
    height: pagePx.value.height + 'px',
    transform: `scale(${zoom.value})`,
    transformOrigin: 'top left',
    position: 'relative' as const,
    background: page.bgColor || '#ffffff',
}));
const pageStyleNoScale = computed(() => ({
    width: pagePx.value.width + 'px',
    height: pagePx.value.height + 'px',
    position: 'relative' as const,
}));
const wrapperStyle = computed(() => ({
    width: pagePx.value.width * zoom.value + 'px',
    height: pagePx.value.height * zoom.value + 'px',
}));
function zoomIn() {
    zoom.value = Math.min(2, Math.round((zoom.value + 0.1) * 10) / 10);
}
function zoomOut() {
    zoom.value = Math.max(0.25, Math.round((zoom.value - 0.1) * 10) / 10);
}

/* Frames state */
const frames = ref<Frame[]>([]);
const selectedId = ref<string | undefined>(undefined);
const selected = computed({
    get: () => frames.value.find((f) => f.id === selectedId.value),
    set: (val) => {
        if (!val) return;
        const idx = frames.value.findIndex((f) => f.id === selectedId.value);
        if (idx >= 0) frames.value[idx] = { ...(frames.value[idx] as any), ...(val as any) };
    },
});
function select(id: string) {
    selectedId.value = id;
}
function updateFrame(updated: Frame) {
    const idx = frames.value.findIndex((f) => f.id === updated.id);
    if (idx >= 0) frames.value[idx] = { ...frames.value[idx], ...updated } as Frame;
}

/* Click-to-insert helper */
function addFrameAt(kind: Kind, xMm: number, yMm: number) {
    const base: BaseFrame = {
        id: crypto.randomUUID(),
        kind,
        name: '',
        xMm,
        yMm,
        wMm: 40,
        hMm: 8,
        strokeColor: '#e03131',
        strokeWidth: 1,
        strokeEnabled: true,
    };
    let f: Frame;
    switch (kind) {
        case 'text':
            f = {
                ...(base as any),
                text: 'Edit via right sidebar',
                fontSize: 12,
                bold: false,
                italic: false,
                underline: false,
                align: 'left',
                rotateDeg: 0,
                style: '',
            };
            break;
        case 'image':
            f = { ...(base as any), src: '/logo.svg', fit: 'contain' } as ImageFrame;
            break;
        case 'repeat':
            f = { ...(base as any), path: 'items' } as RepeatFrame;
            break;
        case 'checkbox':
            f = { ...(base as any), label: 'Checkbox', checked: false } as CheckboxFrame;
            break;
        case 'grid':
            f = { ...(base as any), cols: 8, rows: 1, gapMm: 1 } as GridFrame;
            break;
        case 'line':
            f = { ...(base as any) } as LineFrame;
            break;
        case 'section':
            f = { ...(base as any), title: 'Section' } as SectionFrame;
            break;
        default:
            f = { ...(base as any) } as any;
    }
    frames.value.push(f);
    selectedId.value = f.id;
}

/* Drag & drop onto canvas */
function pagePointFromEvent(e: DragEvent) {
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
    const xPx = (e.clientX - rect.left) / zoom.value;
    const yPx = (e.clientY - rect.top) / zoom.value;
    return { xMm: pxToMm(xPx), yMm: pxToMm(yPx) };
}
function snapMm(v: number) {
    return snap.value ? Math.round(v / gridStepMm.value) * gridStepMm.value : v;
}
// Unified drop handler (palette kind + data leaf)
function onDrop(e: DragEvent) {
    const point = pagePointFromEvent(e);
    const dataLeaf = e.dataTransfer?.getData('text/x-datapath');
    const kindStr = e.dataTransfer?.getData('application/x-kind');
    const snappedX = snapMm(point.xMm),
        snappedY = snapMm(point.yMm);
    if (dataLeaf) {
        frames.value.push({
            id: crypto.randomUUID(),
            kind: 'field',
            name: '',
            xMm: snappedX,
            yMm: snappedY,
            wMm: 40,
            hMm: 8,
            strokeColor: '#e03131',
            strokeWidth: 1,
            path: dataLeaf,
        } as FieldFrame);

        selectedId.value = frames.value[frames.value.length - 1]?.id;
        return;
    }
    if (kindStr) {
        const { kind } = JSON.parse(kindStr) as { kind: Kind };
        addFrameAt(kind, snappedX, snappedY);
        return;
    }
}

/* Delete key support (CHANGE) */
onMounted(() => {
    window.addEventListener('keydown', onKey);
});
onBeforeUnmount(() => {
    window.removeEventListener('keydown', onKey);
});
function onKey(e: KeyboardEvent) {
    const target = e.target as HTMLElement | null;
    // ignore when typing in inputs or contenteditable
    if (target && (['INPUT', 'TEXTAREA'].includes(target.tagName) || target.closest('[contenteditable="true"]'))) return;
    if ((e.key === 'Delete' || e.key === 'Backspace') && selectedId.value) {
        const idx = frames.value.findIndex((f) => f.id === selectedId.value);
        if (idx !== -1) {
            frames.value.splice(idx, 1);
            selectedId.value = undefined;
        }
    }
}

/* Preview helpers */
const sampleData = reactive({
    patient: { name: 'Max Mustermann', birthDate: '1980-05-01' },
    case: { startDate: '2025-08-01', initial: true },
    employer: { name: 'ACME GmbH' },
    items: [{ a: 1 }, { a: 2 }],
});
function fieldLookup(path?: string) {
    if (!path) return '';
    const parts = path.split('.');
    let cur: any = sampleData;
    for (const p of parts) cur = cur?.[p];
    return cur ?? '';
}

const PreviewCell = (props: { frame: Frame; sample: any }) => {
    const f = props.frame as Frame;
    // Only show border if strokeEnabled is not false (default true)
    const border =
        f.kind !== 'line' && f.strokeEnabled !== false && f.strokeWidth
            ? `${f.strokeWidth}px solid ${f.strokeColor || '#e03131'}`
            : undefined;
    const baseStyle: any = {
        width: '100%',
        height: '100%',
        boxSizing: 'border-box',
        background: (f as any).fill || undefined,
        border,
    };
    if (f.kind === 'text') return h('div', { style: { ...baseStyle, padding: '1.5mm', color: '#222' } }, (f as TextFrame).text);
    if (f.kind === 'field')
        return h(
            'div',
            { style: { ...baseStyle, padding: '1.5mm' } },
            String(fieldLookup((f as FieldFrame).path) || (f as FieldFrame).fallback || ''),
        );
    if (f.kind === 'image')
        return h('img', {
            src: (f as ImageFrame).src,
            style: { width: '100%', height: '100%', objectFit: (f as ImageFrame).fit || 'contain' },
        });
    if (f.kind === 'checkbox') {
        const on = (f as CheckboxFrame).path ? !!fieldLookup((f as CheckboxFrame).path) : !!(f as CheckboxFrame).checked;
        return h('div', { style: { ...baseStyle, display: 'flex', alignItems: 'center', gap: '2mm', padding: '1mm' } }, [
            h(
                'div',
                {
                    style: {
                        width: '5mm',
                        height: '5mm',
                        border: '1px solid ' + (f.strokeColor || '#e03131'),
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                    },
                },
                on ? '✓' : '',
            ),
            h('span', (f as CheckboxFrame).label || ''),
        ]);
    }
    if (f.kind === 'grid') {
        const { cols = 8, rows = 1, gapMm = 1 } = f as GridFrame;
        const cells: any[] = [];
        for (let r = 0; r < rows; r++)
            for (let c = 0; c < cols; c++)
                cells.push(h('div', { style: { border: '1px solid ' + (f.strokeColor || '#e03131') } }));
        return h(
            'div',
            { style: { ...baseStyle, display: 'grid', gridTemplateColumns: `repeat(${cols}, 1fr)`, gap: `${gapMm}mm` } },
            cells,
        );
    }
    if (f.kind === 'line')
        return h('div', {
            style: {
                position: 'absolute',
                bottom: '1.5mm',
                left: 0,
                right: 0,
                borderBottom: '1px solid ' + (f.strokeColor || '#e03131'),
            },
        });
    // Removed rotatedText handling
    if (f.kind === 'section') {
        const title = (f as SectionFrame).title || '';
        return h('div', { style: { ...baseStyle, padding: '1.5mm' } }, [
            h('div', { style: { color: '#e03131', fontWeight: 600, marginBottom: '1mm' } }, title),
        ]);
    }
    return h('div', { style: baseStyle }, 'Repeat region');
};

function frameAbsStyle(f: Frame) {
    return { left: mmToPx(f.xMm) + 'px', top: mmToPx(f.yMm) + 'px', width: mmToPx(f.wMm) + 'px', height: mmToPx(f.hMm) + 'px' };
}

/* Presets */
function dropPreset(kind: 'icdRow' | 'stampBox') {
    if (kind === 'icdRow') {
        const baseY = 110,
            colW = 55,
            gap = 5;
        const framesToAdd: Frame[] = Array.from({ length: 4 }).map((_, i) => ({
            id: crypto.randomUUID(),
            kind: 'grid',
            name: `ICD ${i + 1}`,
            xMm: 15 + i * (colW + gap),
            yMm: baseY,
            wMm: colW,
            hMm: 12,
            cols: 6,
            rows: 1,
            gapMm: 1,
            strokeColor: '#e03131',
            strokeWidth: 1,
        })) as any;
        frames.value.push(...framesToAdd);
        return;
    }
    if (kind === 'stampBox') {
        frames.value.push({
            id: crypto.randomUUID(),
            kind: 'section',
            name: 'Stamp',
            xMm: 150,
            yMm: 60,
            wMm: 50,
            hMm: 50,
            title: '',
            strokeColor: '#e03131',
            strokeWidth: 1,
        });
        frames.value.push({
            id: crypto.randomUUID(),
            kind: 'text',
            name: 'StampText',
            xMm: 155,
            yMm: 65,
            wMm: 40,
            hMm: 40,
            text: 'Verbindliches Muster',
            fontSize: 14,
            bold: true,
            rotateDeg: -24,
            style: 'color: #e03131; display: flex; align-items: center; justify-content: center; font-weight: 700;',
        });
    }
}

/* Publish & history */
function publish() {
    const payload = { page: { ...page }, frames: frames.value };
    console.log('Publish payload', payload);
    useToast().add({ title: 'Published (console)', description: 'See devtools for payload JSON.' });
}
const history: string[] = [];
let historyIdx = -1;
function snapshot() {
    history.splice(historyIdx + 1);
    history.push(JSON.stringify({ page, frames: frames.value }));
    historyIdx = history.length - 1;
}
watch([page, frames], () => snapshot(), { deep: true, immediate: true });
const canUndo = computed(() => historyIdx > 0);
const canRedo = computed(() => historyIdx < history.length - 1);
function undo() {
    if (!canUndo.value) return;
    historyIdx--;
    restore(history[historyIdx]!);
}
function redo() {
    if (!canRedo.value) return;
    historyIdx++;
    restore(history[historyIdx]!);
}
function restore(s: string) {
    const st = JSON.parse(s);
    Object.assign(page, st.page);
    frames.value = st.frames;
}

const showPreview = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <!-- Top Toolbar -->
            <UDashboardNavbar title="Layout Editor">
                <template #left>
                    <div class="mx-auto flex items-center gap-2 px-4 py-2">
                        <UButton icon="i-heroicons-arrow-uturn-left" variant="ghost" :disabled="!canUndo" @click="undo" />
                        <UButton icon="i-heroicons-arrow-uturn-right" variant="ghost" :disabled="!canRedo" @click="redo" />
                        <UDivider orientation="vertical" class="mx-1" />
                        <USelect v-model="page.size" :options="pageSizeOptions" class="w-44" />
                        <UButton size="sm" variant="ghost" @click="zoomOut">-</UButton>
                        <div class="w-12 text-center text-sm tabular-nums">{{ Math.round(zoom * 100) }}%</div>
                        <UButton size="sm" variant="ghost" @click="zoomIn">+</UButton>
                        <UDivider orientation="vertical" class="mx-1" />
                        <UToggle v-model="snap" label="Snap" />
                        <USelect v-model="gridStepMm" :options="gridStepOptions" class="w-24" />
                    </div>
                </template>

                <template #right>
                    <UButton icon="i-heroicons-eye" @click="showPreview = true">Preview</UButton>
                    <UButton color="primary" icon="i-heroicons-cloud-arrow-up" @click="publish">Publish</UButton>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <div class="grid grid-cols-8">
                    <!-- Left Sidebar -->
                    <div class="col-span-2 space-y-3 overflow-auto p-3">
                        <UCard>
                            <template #header>Insert</template>
                            <div class="grid grid-cols-1 gap-2 lg:grid-cols-2">
                                <UButton
                                    v-for="item in insertItems"
                                    :key="item.kind"
                                    draggable="true"
                                    variant="soft"
                                    @click="addFrameAt(item.kind, 10, 10)"
                                    >{{ item.label }}</UButton
                                >
                            </div>
                            <p class="mt-2 text-xs text-gray-900 dark:text-white">Tip: click to insert.</p>
                        </UCard>
                        <UCard>
                            <template #header>Components</template>
                            <div class="grid grid-cols-1 gap-2">
                                <UButton variant="soft" @click="dropPreset('icdRow')">ICD-10 Row</UButton>
                                <UButton variant="soft" @click="dropPreset('stampBox')">Stamp Box</UButton>
                            </div>
                        </UCard>
                    </div>

                    <div class="col-span-4 flex h-full w-full">
                        <!-- Center: Canvas -->
                        <section
                            ref="scrollArea"
                            class="relative flex-1 overflow-auto bg-[linear-gradient(45deg,rgba(0,0,0,0.02)_25%,transparent_25%),linear-gradient(-45deg,rgba(0,0,0,0.02)_25%,transparent_25%),linear-gradient(45deg,transparent_75%,rgba(0,0,0,0.02)_75%),linear-gradient(-45deg,transparent_75%,rgba(0,0,0,0.02)_75%)] bg-[length:20px_20px,20px_20px,20px_20px,20px_20px] bg-[position:0_0,0_10px,10px_-10px,-10px_0]"
                        >
                            <!-- Rulers -->
                            <div class="sticky left-0 top-0 z-10 flex bg-transparent backdrop-blur">
                                <div class="h-6 w-6 border-b border-r bg-gray-300/10" />
                                <div class="relative h-6 flex-1 overflow-hidden bg-gray-300/10">
                                    <Ruler :length-px="pagePx.width" orientation="horizontal" :zoom="zoom" />
                                </div>
                            </div>

                            <div class="absolute bottom-0 left-0 top-6 z-10 w-6 bg-gray-300/10">
                                <Ruler :length-px="pagePx.height" orientation="vertical" :zoom="zoom" />
                            </div>

                            <!-- Drop target wrapper (centers page) -->
                            <div class="flex items-start justify-center pl-6 pt-8">
                                <div class="relative" :style="wrapperStyle" @dragover.prevent @drop.prevent="onDrop">
                                    <!-- Page -->
                                    <div class="page border bg-white shadow" :style="pageStyle">
                                        <!-- Background color -->
                                        <div
                                            class="pointer-events-none absolute inset-0 h-full w-full select-none"
                                            :style="{ background: page.bgColor, zIndex: 0 }"
                                        />
                                        <!-- Background behind -->
                                        <img
                                            v-if="page.backgroundUrl && page.bgLayer === 'behind'"
                                            :src="page.backgroundUrl"
                                            class="pointer-events-none absolute inset-0 h-full w-full select-none"
                                            :style="{ opacity: page.bgOpacity, zIndex: 1 }"
                                        />

                                        <!-- Grid overlay -->
                                        <svg
                                            class="pointer-events-none absolute inset-0"
                                            :width="pagePx.width"
                                            :height="pagePx.height"
                                        >
                                            <line
                                                v-for="x in Math.floor(page.widthMm / gridStepMm)"
                                                :key="'vx-' + x"
                                                :x1="mmToPx(x * gridStepMm)"
                                                y1="0"
                                                :x2="mmToPx(x * gridStepMm)"
                                                :y2="pagePx.height"
                                                stroke="rgba(0,0,0,0.05)"
                                                stroke-width="1"
                                            />
                                            <line
                                                v-for="y in Math.floor(page.heightMm / gridStepMm)"
                                                :key="'hz-' + y"
                                                x1="0"
                                                :y1="mmToPx(y * gridStepMm)"
                                                :x2="pagePx.width"
                                                :y2="mmToPx(y * gridStepMm)"
                                                stroke="rgba(0,0,0,0.05)"
                                                stroke-width="1"
                                            />
                                        </svg>

                                        <!-- Frames -->
                                        <FrameNode
                                            v-for="f in frames"
                                            :key="f.id"
                                            :frame="f"
                                            :selected="selectedId === f.id"
                                            :zoom="zoom"
                                            :snap="snap"
                                            :grid-step-mm="gridStepMm"
                                            @select="select(f.id)"
                                            @update:frame="updateFrame"
                                        />

                                        <!-- Background overlay on top -->
                                        <img
                                            v-if="page.backgroundUrl && page.bgLayer === 'overlay'"
                                            :src="page.backgroundUrl"
                                            class="pointer-events-none absolute inset-0 h-full w-full select-none"
                                            :style="{ opacity: page.bgOpacity, mixBlendMode: page.bgBlend }"
                                        />
                                    </div>
                                </div>
                            </div>
                        </section>
                    </div>

                    <!-- Right: Properties -->
                    <div class="col-span-2 space-y-3 overflow-auto p-3">
                        <UCard>
                            <template #header>Page</template>
                            <div class="grid grid-cols-2 gap-2">
                                <UFormGroup label="Width (mm)">
                                    <UInput v-model.number="page.widthMm" type="number" step="1" />
                                </UFormGroup>
                                <UFormGroup label="Height (mm)">
                                    <UInput v-model.number="page.heightMm" type="number" step="1" />
                                </UFormGroup>
                                <UFormGroup label="Background" class="col-span-2">
                                    <div class="flex items-center gap-2">
                                        <UButton size="xs" icon="i-heroicons-photo" @click="pickBg">Choose</UButton>
                                        <UToggle v-model="page.bgLocked" label="Lock" />
                                        <ColorPickerClient v-model="page.bgColor" class="ml-2" />
                                    </div>
                                    <input ref="bgInput" type="file" accept="image/*" class="hidden" @change="onBgChange" />
                                </UFormGroup>
                                <UFormGroup label="Opacity"
                                    ><URange v-model="page.bgOpacity" :min="0" :max="1" :step="0.05" />
                                </UFormGroup>
                                <UFormGroup label="Layer"
                                    ><USelect v-model="page.bgLayer" :options="bgLayerOptions" />
                                </UFormGroup>
                                <UFormGroup label="Blend"
                                    ><USelect v-model="page.bgBlend" :options="bgBlendOptions" />
                                </UFormGroup>
                            </div>
                        </UCard>

                        <UCard>
                            <template #header>Properties</template>
                            <div v-if="selected" class="space-y-3">
                                <UInput v-model="selected.name" placeholder="Name" />
                                <div class="grid grid-cols-2 gap-2">
                                    <UFormGroup label="X (mm)">
                                        <UInput v-model.number="selected.xMm" type="number" step="0.5" />
                                    </UFormGroup>
                                    <UFormGroup label="Y (mm)">
                                        <UInput v-model.number="selected.yMm" type="number" step="0.5" />
                                    </UFormGroup>
                                    <UFormGroup label="W (mm)">
                                        <UInput v-model.number="selected.wMm" type="number" step="0.5" />
                                    </UFormGroup>
                                    <UFormGroup label="H (mm)">
                                        <UInput v-model.number="selected.hMm" type="number" step="0.5" />
                                    </UFormGroup>
                                </div>
                                <USelect v-model="selected.kind" :options="kindOptions" />
                                <UDivider />
                                <div class="grid grid-cols-2 gap-2">
                                    <UFormGroup label="Stroke">
                                        <ColorPickerClient v-model="selected.strokeColor" />
                                    </UFormGroup>
                                    <UFormGroup label="Width">
                                        <UInput v-model.number="selected.strokeWidth" type="number" step="1" />
                                    </UFormGroup>
                                    <UFormGroup label="Show Stroke" class="col-span-2">
                                        <UToggle v-model="selected.strokeEnabled" label="Show border (stroke)" />
                                    </UFormGroup>
                                    <UFormGroup label="Fill" class="col-span-2">
                                        <UInput v-model="selected.fill" placeholder="transparent" />
                                    </UFormGroup>
                                </div>
                                <div v-if="selected.kind === 'field'">
                                    <UFormGroup label="Path">
                                        <UInput v-model="(selected as any).path" placeholder="patient.name" />
                                    </UFormGroup>
                                    <UFormGroup label="Fallback">
                                        <UInput v-model="(selected as any).fallback" placeholder="—" />
                                    </UFormGroup>
                                </div>
                                <div v-if="selected.kind === 'text'">
                                    <UFormGroup label="Text"
                                        ><UTextarea v-model="(selected as any).text" :rows="5" />
                                    </UFormGroup>
                                </div>
                                <div v-if="selected.kind === 'image'">
                                    <UFormGroup label="Image URL">
                                        <UInput v-model="(selected as any).src" placeholder="/logo.png" />
                                    </UFormGroup>
                                    <USelect v-model="(selected as any).fit" :options="imageFitOptions" />
                                </div>
                                <div v-if="selected.kind === 'checkbox'">
                                    <UFormGroup label="Path">
                                        <UInput v-model="(selected as any).path" placeholder="case.initial" />
                                    </UFormGroup>
                                    <UFormGroup label="Label">
                                        <UInput v-model="(selected as any).label" placeholder="Erstbescheinigung" />
                                    </UFormGroup>
                                    <UFormGroup label="Checked (preview)">
                                        <UToggle v-model="(selected as any).checked" label="Checked (preview)" />
                                    </UFormGroup>
                                </div>
                                <div v-if="selected.kind === 'grid'">
                                    <div class="grid grid-cols-3 gap-2">
                                        <UFormGroup label="Cols">
                                            <UInput v-model.number="(selected as any).cols" type="number" />
                                        </UFormGroup>
                                        <UFormGroup label="Rows">
                                            <UInput v-model.number="(selected as any).rows" type="number" />
                                        </UFormGroup>
                                        <UFormGroup label="Gap (mm)">
                                            <UInput v-model.number="(selected as any).gapMm" type="number" step="0.5" />
                                        </UFormGroup>
                                    </div>
                                </div>
                                <div v-if="selected.kind === 'text'">
                                    <UFormGroup label="Text">
                                        <UTextarea v-model="(selected as any).text" :rows="5" />
                                    </UFormGroup>
                                    <div class="grid grid-cols-2 gap-2">
                                        <UFormGroup label="Font Size (pt)">
                                            <UInput
                                                v-model.number="(selected as any).fontSize"
                                                type="number"
                                                min="6"
                                                max="72"
                                                step="1"
                                            />
                                        </UFormGroup>
                                        <UFormGroup label="Align">
                                            <USelect
                                                v-model="(selected as any).align"
                                                :options="[
                                                    { label: 'Left', value: 'left' },
                                                    { label: 'Center', value: 'center' },
                                                    { label: 'Right', value: 'right' },
                                                ]"
                                            />
                                        </UFormGroup>
                                    </div>
                                    <div class="flex gap-2">
                                        <UToggle v-model="(selected as any).bold" label="Bold" />
                                        <UToggle v-model="(selected as any).italic" label="Italic" />
                                        <UToggle v-model="(selected as any).underline" label="Underline" />
                                    </div>
                                    <UFormGroup label="Rotation (deg)">
                                        <UInput
                                            v-model.number="(selected as any).rotateDeg"
                                            type="number"
                                            min="-180"
                                            max="180"
                                            step="1"
                                        />
                                    </UFormGroup>
                                    <UFormGroup label="Custom Style (CSS)">
                                        <UInput
                                            v-model="(selected as any).style"
                                            placeholder="color: red; background: yellow;"
                                        />
                                    </UFormGroup>
                                </div>
                                <div v-if="selected.kind === 'section'">
                                    <UFormGroup label="Title">
                                        <UInput v-model="(selected as any).title" placeholder="AU-begründende Diagnose(n)" />
                                    </UFormGroup>
                                </div>
                            </div>
                            <div v-else class="text-sm text-gray-900 dark:text-white">
                                Select a frame to edit its properties.
                            </div>
                        </UCard>
                    </div>
                </div>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>

    <!-- Preview Modal -->
    <UModal v-model="showPreview" fullscreen>
        <UCard
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                base: 'flex flex-1 flex-col',
                body: { base: 'flex flex-1 flex-col' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">Preview (sample data)</h3>

                    <UButton
                        class="-my-1"
                        color="gray"
                        variant="ghost"
                        icon="i-mdi-window-close"
                        @click="showPreview = false"
                    />
                </div>
            </template>

            <div class="flex flex-1 items-center justify-center">
                <div class="overflow-auto p-4">
                    <div class="mx-auto" :style="{ width: pagePx.width + 'px' }">
                        <div class="page border" :style="pageStyleNoScale">
                            <div
                                class="pointer-events-none absolute inset-0 h-full w-full select-none"
                                :style="{ background: page.bgColor, zIndex: 0 }"
                            />
                            <img
                                v-if="page.backgroundUrl && page.bgLayer !== 'behind'"
                                :src="page.backgroundUrl"
                                class="pointer-events-none absolute inset-0 h-full w-full select-none"
                                :style="{ opacity: page.bgOpacity, mixBlendMode: page.bgBlend, zIndex: 1 }"
                            />
                            <img
                                v-if="page.backgroundUrl && page.bgLayer === 'behind'"
                                :src="page.backgroundUrl"
                                class="pointer-events-none absolute inset-0 h-full w-full select-none"
                                :style="{ opacity: page.bgOpacity, zIndex: 1 }"
                            />
                            <div v-for="f in frames" :key="'pv-' + f.id" class="absolute" :style="frameAbsStyle(f)">
                                <component :is="PreviewCell" :frame="f" :sample="sampleData" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" @click="showPreview = false">Close</UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

<style scoped>
.page {
    width: 210mm;
    height: 297mm;
}
</style>

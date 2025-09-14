<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import { useToast } from '#imports';
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import FrameNode from './FrameNode.vue';
import PreviewCell from './PreviewCell.vue';
import Ruler from './Ruler.vue';
import type { Frame, WidgetFrame } from './types';

const { fileUpload } = useAppConfig();

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

function onBgChange(file: File | null | undefined) {
    if (!file) return;
    page.backgroundUrl = URL.createObjectURL(file);
}

const formData = reactive<Record<string, unknown>>({});

function setData(path: string, v: unknown) {
    formData[path] = v;
}

/* Canvas */
const zoom = ref(0.9);
const gridStepMm = ref(2);
const gridStepOptions = [
    { label: '1 mm', value: 1 },
    { label: '2 mm', value: 2 },
    { label: '5 mm', value: 5 },
];

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
    const pt = pagePointFromEvent(e);
    const x = snapMm(pt.xMm),
        y = snapMm(pt.yMm);

    // widgets only
    const widgetStr = e.dataTransfer?.getData('application/x-widget');
    if (widgetStr) {
        addWidgetAt(JSON.parse(widgetStr) as NewWidgetPayload, x, y);
        return;
    }

    // If you still support non-input shapes (lines, images), keep that branch.
    // Anything that used to create 'field'/'text' input frames should be removed.
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

function frameAbsStyle(f: Frame) {
    return { left: mmToPx(f.xMm) + 'px', top: mmToPx(f.yMm) + 'px', width: mmToPx(f.wMm) + 'px', height: mmToPx(f.hMm) + 'px' };
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
    const st = JSON.parse(s) as { page: typeof page; frames: Frame[] };
    Object.assign(page, st.page);
    frames.value = st.frames;
}

type NewWidgetPayload = { widget: 'text' | 'textarea' | 'boxed-text'; rows?: number; boxCharCount?: number };

function onWidgetDragStart(e: DragEvent, payload: NewWidgetPayload) {
    e.dataTransfer?.setData('application/x-widget', JSON.stringify(payload));
    e.dataTransfer?.setData('text/plain', 'widget'); // cursor hint
}

function addWidgetAt(payload: NewWidgetPayload, xMm: number, yMm: number) {
    const baseW = payload.widget === 'textarea' ? 60 : payload.widget === 'boxed-text' ? 40 : 40;
    const baseH = payload.widget === 'textarea' ? 20 : 8;

    const f: WidgetFrame = {
        id: crypto.randomUUID(),
        kind: 'widget',
        name: '',
        xMm,
        yMm,
        wMm: baseW,
        hMm: baseH,
        strokeColor: '#000000',
        strokeWidth: 1,
        strokeEnabled: true,
        widget: payload.widget,
        rows: payload.rows,
        boxCharCount: payload.boxCharCount,
    };
    frames.value.push(f as Frame);
    selectedId.value = f.id;
}

const showPreview = ref(false);
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <!-- Top Toolbar -->
            <UDashboardNavbar title="Layout Editor">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #left>
                    <UDashboardSidebarCollapse />

                    <div class="mx-auto flex items-center gap-2 px-4 py-2">
                        <UButton icon="i-heroicons-arrow-uturn-left" variant="ghost" :disabled="!canUndo" @click="undo" />
                        <UButton icon="i-heroicons-arrow-uturn-right" variant="ghost" :disabled="!canRedo" @click="redo" />
                        <USeparator orientation="vertical" class="mx-1" />
                        <USelect v-model="page.size" :items="pageSizeOptions" class="w-44" />
                        <UButton size="sm" variant="ghost" @click="zoomOut">-</UButton>
                        <div class="w-12 text-center text-sm tabular-nums">{{ Math.round(zoom * 100) }}%</div>
                        <UButton size="sm" variant="ghost" @click="zoomIn">+</UButton>
                        <USeparator orientation="vertical" class="mx-1" />
                        <USwitch v-model="snap" label="Snap" />
                        <USelect v-model="gridStepMm" :items="gridStepOptions" class="w-24" />
                    </div>
                </template>

                <template #right>
                    <UButton icon="i-mdi-eye" @click="showPreview = true">Preview</UButton>
                    <UButton color="primary" icon="i-heroicons-cloud-arrow-up" @click="publish">Publish</UButton>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="grid h-full grid-cols-8">
                <!-- Left Sidebar -->
                <div class="col-span-2 space-y-3 overflow-auto p-2">
                    <!-- In <template> Left Sidebar, below the existing <UCard> Insert -->
                    <UCard>
                        <template #header>Widgets</template>
                        <div class="grid grid-cols-1 gap-2 lg:grid-cols-2">
                            <UButton
                                variant="soft"
                                draggable="true"
                                @dragstart="onWidgetDragStart($event, { widget: 'text' })"
                                @click="addWidgetAt({ widget: 'text' }, 10, 10)"
                            >
                                Text
                            </UButton>

                            <UButton
                                variant="soft"
                                draggable="true"
                                @dragstart="onWidgetDragStart($event, { widget: 'textarea', rows: 3 })"
                                @click="addWidgetAt({ widget: 'textarea', rows: 3 }, 10, 10)"
                            >
                                Textarea
                            </UButton>

                            <UButton
                                variant="soft"
                                draggable="true"
                                @dragstart="onWidgetDragStart($event, { widget: 'boxed-text', boxCharCount: 8 })"
                                @click="addWidgetAt({ widget: 'boxed-text', boxCharCount: 8 }, 10, 10)"
                            >
                                Boxed text
                            </UButton>
                        </div>
                        <p class="mt-2 text-xs text-highlighted">Drag onto the page or click to insert.</p>
                    </UCard>
                </div>

                <div class="col-span-4 flex h-full w-full">
                    <!-- Center: Canvas -->
                    <section
                        ref="scrollArea"
                        class="relative flex-1 overflow-auto bg-[linear-gradient(45deg,rgba(0,0,0,0.02)_25%,transparent_25%),linear-gradient(-45deg,rgba(0,0,0,0.02)_25%,transparent_25%),linear-gradient(45deg,transparent_75%,rgba(0,0,0,0.02)_75%),linear-gradient(-45deg,transparent_75%,rgba(0,0,0,0.02)_75%)] bg-size-[20px_20px,20px_20px,20px_20px,20px_20px] bg-position-[0_0,0_10px,10px_-10px,-10px_0]"
                    >
                        <!-- Rulers -->
                        <div class="sticky top-0 left-0 z-10 flex bg-transparent backdrop-blur-sm">
                            <div class="h-6 w-6 border-r border-b bg-neutral-300/10" />
                            <div class="relative h-6 flex-1 overflow-hidden bg-neutral-300/10">
                                <Ruler :length-px="pagePx.width" orientation="horizontal" :zoom="zoom" />
                            </div>
                        </div>

                        <div class="absolute top-6 bottom-0 left-0 z-10 w-6 bg-neutral-300/10">
                            <Ruler :length-px="pagePx.height" orientation="vertical" :zoom="zoom" />
                        </div>

                        <!-- Drop target wrapper (centers page) -->
                        <div class="flex items-start justify-center pt-8 pl-6">
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
                                    <template v-for="f in frames" :key="f.id">
                                        <FrameNode
                                            v-if="f.kind === 'widget'"
                                            :node="f as WidgetFrame"
                                            mode="fill"
                                            :data="formData"
                                            :set-data="setData"
                                            @select="select(f.id)"
                                            @update:frame="updateFrame"
                                        />
                                    </template>

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
                <div class="col-span-2 space-y-3 overflow-auto p-2">
                    <UCard>
                        <template #header>Page</template>
                        <div class="grid grid-cols-2 gap-2">
                            <UFormField label="Width (mm)">
                                <UInput v-model.number="page.widthMm" type="number" step="1" />
                            </UFormField>
                            <UFormField label="Height (mm)">
                                <UInput v-model.number="page.heightMm" type="number" step="1" />
                            </UFormField>
                            <UFormField label="Background" class="col-span-2">
                                <div class="flex items-center gap-2">
                                    <USwitch v-model="page.bgLocked" label="Lock" />
                                    <ColorPicker v-model="page.bgColor" class="ml-2" />
                                </div>

                                <UFileUpload
                                    :accept="fileUpload.types.images.join(',')"
                                    :placeholder="$t('common.image')"
                                    :label="$t('common.file_upload_label')"
                                    :description="$t('common.allowed_file_types')"
                                    @update:model-value="onBgChange"
                                />
                            </UFormField>
                            <UFormField label="Opacity"
                                ><USlider v-model="page.bgOpacity" :min="0" :max="1" :step="0.05" />
                            </UFormField>
                            <UFormField label="Layer"><USelect v-model="page.bgLayer" :items="bgLayerOptions" /> </UFormField>
                            <UFormField label="Blend"><USelect v-model="page.bgBlend" :items="bgBlendOptions" /> </UFormField>
                        </div>
                    </UCard>

                    <UCard>
                        <template #header>Properties</template>
                        <div v-if="selected" class="space-y-3">
                            <UInput v-model="selected.name" placeholder="Name" />
                            <div class="grid grid-cols-2 gap-2">
                                <UFormField label="X (mm)">
                                    <UInput v-model.number="selected.xMm" type="number" step="0.5" />
                                </UFormField>
                                <UFormField label="Y (mm)">
                                    <UInput v-model.number="selected.yMm" type="number" step="0.5" />
                                </UFormField>
                                <UFormField label="W (mm)">
                                    <UInput v-model.number="selected.wMm" type="number" step="0.5" />
                                </UFormField>
                                <UFormField label="H (mm)">
                                    <UInput v-model.number="selected.hMm" type="number" step="0.5" />
                                </UFormField>
                            </div>
                            <USeparator />
                            <div class="grid grid-cols-2 gap-2">
                                <UFormField label="Stroke">
                                    <ColorPicker v-model="selected.strokeColor" />
                                </UFormField>
                                <UFormField label="Width">
                                    <UInput v-model.number="selected.strokeWidth" type="number" step="1" />
                                </UFormField>
                                <UFormField label="Show Stroke" class="col-span-2">
                                    <USwitch v-model="selected.strokeEnabled" label="Show border (stroke)" />
                                </UFormField>
                                <UFormField label="Fill" class="col-span-2">
                                    <UInput v-model="selected.fill" placeholder="transparent" />
                                </UFormField>
                            </div>
                            <!-- Only when a widget is selected -->
                            <div v-if="selected?.kind === 'widget'">
                                <UFormGroup label="Widget type">
                                    <USelect
                                        v-model="(selected as any).widget"
                                        :options="[
                                            { label: 'Text', value: 'text' },
                                            { label: 'Textarea', value: 'textarea' },
                                            { label: 'Boxed text', value: 'boxed-text' },
                                        ]"
                                    />
                                </UFormGroup>

                                <UFormGroup label="Bind to field (path)">
                                    <PathPicker v-model="(selected as any).binding" />
                                    <!-- binding.path, formatter, fit -->
                                </UFormGroup>

                                <UFormGroup v-if="(selected as any).widget === 'textarea'" label="Rows">
                                    <UInput v-model.number="(selected as any).rows" type="number" :min="2" :max="20" />
                                </UFormGroup>

                                <UFormGroup v-if="(selected as any).widget === 'boxed-text'" label="Box count">
                                    <UInput v-model.number="(selected as any).boxCharCount" type="number" :min="1" :max="64" />
                                </UFormGroup>

                                <UFormGroup label="Max characters">
                                    <UInput v-model.number="(selected as any).maxChars" type="number" :min="1" :max="512" />
                                </UFormGroup>

                                <UFormGroup label="Fit strategy">
                                    <USelect
                                        v-model="(selected as any).binding.fit"
                                        :options="['wrap', 'truncate', 'shrink']"
                                    />
                                </UFormGroup>

                                <UFormGroup v-if="(selected as any).binding?.fit === 'shrink'" label="Min font size">
                                    <UInput
                                        v-model.number="(selected as any).binding.minFontSize"
                                        type="number"
                                        :min="6"
                                        :max="18"
                                    />
                                </UFormGroup>
                            </div>
                        </div>
                        <div v-else class="text-sm text-highlighted">Select a frame to edit its properties.</div>
                    </UCard>
                </div>
            </div>
        </template>
    </UDashboardPanel>

    <!-- Preview Modal -->
    <UModal v-model:open="showPreview" title="Preview (sample data)" fullscreen>
        <template #body>
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
                            <div v-for="frame in frames" :key="'pv-' + frame.id" class="absolute" :style="frameAbsStyle(frame)">
                                <PreviewCell v-if="frame.kind === 'widget'" :node="frame" :data="formData" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" :label="$t('common.close')" @click="showPreview = false" />
            </UButtonGroup>
        </template>
    </UModal>
</template>

<style scoped>
.page {
    width: 210mm;
    height: 297mm;
}
</style>

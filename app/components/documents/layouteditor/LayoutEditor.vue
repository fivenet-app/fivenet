<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import SVGCanvas, { type SvgCanvasAPI } from './SVGCanvas.vue';

const sc = useTemplateRef<typeof SVGCanvas>('sc');
const docTitle = ref('Untitled.svg');
const phKey = ref('');
const gridOn = ref(true);
const snapOn = ref(true);
const paperSize = ref<'A4 Portrait' | 'A4 Landscape' | 'Letter' | 'Custom'>('A4 Portrait');
const canvasDims = reactive<[number, number]>([1000, 700]);

async function newDoc() {
    // reset by setting a minimal SVG; fallback to clearing root if your build exposes it
    await sc.value?.setSvg?.(`<svg width="${canvasDims[0]}" height="${canvasDims[1]}" xmlns="http://www.w3.org/2000/svg"/>`);
}
async function openSvg() {
    const file = await pickFile('.svg');
    if (!file) return;
    const text = await file.text();
    (await sc.value?.loadSvg?.(text)) || sc.value?.setSvg?.(text);
}
async function saveSvg() {
    const svg = await sc.value?.getSvgString?.();
    if (!svg) return;
    downloadText(svg, sanitizeName(docTitle.value || 'Untitled') + '.svg', 'image/svg+xml');
}
async function exportPng() {
    const svg = await sc.value?.getSvgString?.();
    if (!svg) return;
    const png = await rasterizeSvgToPng(svg, canvasDims[0], canvasDims[1]);
    downloadBlob(png, sanitizeName(docTitle.value || 'Untitled') + '.png');
}
function zoomIn() {
    sc.value?.zoom(1.1);
}
function zoomOut() {
    sc.value?.zoom(1 / 1.1);
}
function fitToContent() {
    sc.value?.fit();
}
watch(snapOn, () => {
    sc.value?.toggleSnap();
});

function addText() {
    sc.value?.addText('Sample text', 120, 120, 18);
}
function addRect() {
    sc.value?.addRect(100, 180, 160, 80);
}
function addLine() {
    sc.value?.addRaw({
        element: 'line',
        attr: { x1: 100, y1: 300, x2: 320, y2: 300, stroke: '#000', 'stroke-width': 1 },
    });
}

// Sidebar actions
function insertPlaceholder() {
    if (!phKey.value) return;
    sc.value?.insertPlaceholder?.(phKey.value, 100, 150, 240);
    phKey.value = '';
}

const paperOptions = ['A4 Portrait', 'A4 Landscape', 'Letter', 'Custom'];
function applyPaper() {
    const map: Record<string, [number, number]> = {
        'A4 Portrait': [794, 1123], // 96dpi approx (8.27in x 11.69in)
        'A4 Landscape': [1123, 794],
        Letter: [816, 1056],
    };
    const dims = map[paperSize.value] || canvasDims;
    canvasDims[0] = dims[0];
    canvasDims[1] = dims[1];
}

// Grid background (pure CSS so it doesn’t touch the SVG)
const gridSize = 16;
const gridStyle = computed(() => ({
    backgroundImage: `
    linear-gradient(to right, rgba(0,0,0,.08) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(0,0,0,.08) 1px, transparent 1px)`,
    backgroundSize: `${gridSize}px ${gridSize}px, ${gridSize}px ${gridSize}px`,
}));

// Tool menu (example)
const toolItems = [
    [
        {
            label: 'Select',
            icon: 'i-heroicons-cursor-arrow-rays',
            click: () => sc.value?.setMode?.('select'),
        },
        {
            label: 'Pan',
            icon: 'i-heroicons-hand-raised',
            click: () => sc.value?.setMode?.('pan'),
        },
    ],
    [
        {
            label: 'Rectangle',
            icon: 'i-heroicons-square-2-stack',
            click: () => addRect(),
        },
        {
            label: 'Line',
            icon: 'i-heroicons-minus',
            click: () => addLine(),
        },
        {
            label: 'Text',
            icon: 'i-heroicons-chat-bubble-left-right',
            click: () => addText(),
        },
    ],
];

// Connect lifecycle from the mounted canvas
function onReady(api: SvgCanvasAPI) {
    console.log('svgcanvas ready', api);
}
function onError(e: unknown) {
    console.error('svgcanvas init error', e);
}

// File helpers
async function pickFile(accept: string) {
    return new Promise<File | null>((resolve) => {
        const inp = document.createElement('input');
        inp.type = 'file';
        inp.accept = accept;
        inp.onchange = () => resolve(inp.files?.[0] ?? null);
        inp.click();
    });
}
function sanitizeName(s: string) {
    return s.replace(/[^\w\-.]+/g, '_');
}
function downloadText(text: string, name: string, mime = 'text/plain') {
    const blob = new Blob([text], { type: mime });
    downloadBlob(blob, name);
}
function downloadBlob(blob: Blob, name: string) {
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = name;
    a.click();
    URL.revokeObjectURL(url);
}

// Quick client-side rasterizer; good enough for forms/preview
async function rasterizeSvgToPng(svg: string, w: number, h: number) {
    const svgBlob = new Blob([svg], { type: 'image/svg+xml' });
    const url = URL.createObjectURL(svgBlob);
    const img = new Image();
    img.decoding = 'async';
    img.src = url;
    await img.decode();
    const canvas = document.createElement('canvas');
    canvas.width = w;
    canvas.height = h;
    const ctx = canvas.getContext('2d')!;
    ctx.drawImage(img, 0, 0, w, h);
    URL.revokeObjectURL(url);
    const blob: Blob = await new Promise((res) => canvas.toBlob((b) => res(b!), 'image/png'));
    return blob;
}
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar title="Layout Editor">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #left>
                    <div class="flex items-center gap-2">
                        <UButton size="sm" icon="i-heroicons-document-plus" @click="newDoc">New</UButton>
                        <UButton size="sm" icon="i-heroicons-folder-open" @click="openSvg">Open</UButton>
                        <UButton size="sm" icon="i-heroicons-arrow-down-tray" @click="saveSvg">Save</UButton>
                        <UButton size="sm" icon="i-heroicons-photo" @click="exportPng">Export PNG</UButton>
                        <USeparator orientation="vertical" class="mx-2" />
                        <UButtonGroup size="sm">
                            <UButton icon="i-heroicons-magnifying-glass-plus" @click="zoomIn" />
                            <UButton icon="i-heroicons-magnifying-glass-minus" @click="zoomOut" />
                            <UButton icon="i-heroicons-arrows-pointing-out" @click="fitToContent">Fit</UButton>
                        </UButtonGroup>
                        <USwitch
                            v-model="gridOn"
                            checked-icon="i-heroicons-squares-2x2"
                            unchecked-icon="i-heroicons-squares-2x2"
                            class="ml-3"
                        />
                        <span class="ml-1 text-xs opacity-70">Grid</span>
                        <USwitch
                            v-model="snapOn"
                            checked-icon="i-heroicons-bolt"
                            unchecked-icon="i-heroicons-bolt"
                            class="ml-3"
                        />
                        <span class="ml-1 text-xs opacity-70">Snap</span>
                    </div>
                </template>

                <template #right>
                    <UInput v-model="docTitle" size="sm" placeholder="Untitled.svg" class="w-48" />
                    <UDropdownMenu :items="toolItems" :popper="{ placement: 'bottom-end' }">
                        <UButton size="sm" icon="i-heroicons-wrench-screwdriver" class="ml-2">Tools</UButton>
                    </UDropdownMenu>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <!-- Body: sidebar + canvas -->
            <div class="grid h-full grid-cols-[280px_1fr] overflow-hidden">
                <!-- Sidebar -->
                <aside class="space-y-4 overflow-y-auto border-r p-4">
                    <UCard>
                        <template #header><div class="font-medium">Placeholders</div></template>
                        <div class="flex items-center gap-2">
                            <UInput v-model="phKey" placeholder="patient_name" size="sm" />
                            <UButton size="sm" icon="i-heroicons-plus" :disabled="!phKey" @click="insertPlaceholder" />
                        </div>
                        <div class="mt-2 text-xs opacity-60">
                            Click + to insert <code>{{ phKey || 'field_key' }}</code> with an underline.
                        </div>
                    </UCard>

                    <UCard>
                        <template #header><div class="font-medium">Canvas</div></template>
                        <div class="grid gap-2">
                            <USelectMenu v-model="paperSize" :items="paperOptions" size="sm" />
                            <UButton size="sm" @click="applyPaper">Apply Paper Size</UButton>
                        </div>
                    </UCard>

                    <UCard>
                        <template #header><div class="font-medium">Add</div></template>
                        <div class="grid grid-cols-2 gap-2">
                            <UButton size="sm" @click="addText">Text</UButton>
                            <UButton size="sm" @click="addRect">Rect</UButton>
                            <UButton size="sm" @click="addLine">Line</UButton>
                        </div>
                    </UCard>
                </aside>

                <!-- Canvas area -->
                <div class="relative flex-1 overflow-y-hidden">
                    <!-- Grid backplate -->
                    <div v-show="gridOn" class="pointer-events-none absolute inset-0" :style="gridStyle" />
                    <!-- The actual svgcanvas -->
                    <SVGCanvas
                        ref="sc"
                        :config="{ dimensions: canvasDims, gridSnapping: snapOn }"
                        @ready="onReady"
                        @error="onError"
                    />
                </div>
            </div>
        </template>
    </UDashboardPanel>
</template>

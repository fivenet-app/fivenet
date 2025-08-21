<script setup lang="ts">
import type { Frame } from './types';

const props = withDefaults(
    defineProps<{
        selected?: boolean;
        zoom: number;
        snap: boolean;
        gridStepMm?: number;
    }>(),
    {
        selected: false,
        gridStepMm: 2,
    },
);

const frame = defineModel<Frame>('frame', { required: true });

const emit = defineEmits<{
    (e: 'select'): void;
}>();

const dragging = ref(false);
const resizing = ref(false);
const start = ref({
    mx: 0,
    my: 0,
    xMm: 0,
    yMm: 0,
    wMm: 0,
    hMm: 0,
});

const mmToPx = (mm: number) => mm * 3.7795275591;
const pxToMm = (px: number) => px / 3.7795275591;

const snapMmLocal = (v: number) => (props.snap ? Math.round(v / props.gridStepMm) * props.gridStepMm : v);

const styleRect = computed(() => ({
    left: mmToPx(frame.value.xMm) + 'px',
    top: mmToPx(frame.value.yMm) + 'px',
    width: mmToPx(frame.value.wMm) + 'px',
    height: mmToPx(frame.value.hMm) + 'px',
}));

const containerClass = computed(() => [
    'absolute group',
    frame.value.strokeWidth ? 'ring-0 border' : 'ring-1 ring-gray-300',
    props.selected ? 'ring-2 ring-primary-500' : '',
]);

const containerStyle = computed(() => ({
    ...styleRect.value,
    border: frame.value.strokeWidth ? `${frame.value.strokeWidth}px solid ${frame.value.strokeColor || '#e03131'}` : undefined,
    background: frame.value.fill || undefined,
}));

const onDown = (e: PointerEvent) => {
    dragging.value = true;
    const target = e.target as Element | null;
    if (target && typeof (target as HTMLElement).setPointerCapture === 'function') {
        (target as HTMLElement).setPointerCapture(e.pointerId);
    }
    start.value = {
        mx: e.clientX,
        my: e.clientY,
        xMm: frame.value.xMm,
        yMm: frame.value.yMm,
        wMm: frame.value.wMm,
        hMm: frame.value.hMm,
    };
    emit('select');
};

const onMove = (e: PointerEvent) => {
    if (!dragging.value && !resizing.value) return;
    const dxPx = (e.clientX - start.value.mx) / props.zoom;
    const dyPx = (e.clientY - start.value.my) / props.zoom;
    if (resizing.value) {
        // Only update width and height when resizing
        const w = snapMmLocal(start.value.wMm + pxToMm(dxPx));
        const h = snapMmLocal(start.value.hMm + pxToMm(dyPx));
        frame.value = {
            ...frame.value,
            wMm: Math.max(5, w),
            hMm: Math.max(5, h),
        };
    } else if (dragging.value) {
        // Only update position when dragging
        const x = snapMmLocal(start.value.xMm + pxToMm(dxPx));
        const y = snapMmLocal(start.value.yMm + pxToMm(dyPx));
        frame.value = {
            ...frame.value,
            xMm: x,
            yMm: y,
        };
    }
};

const onUp = () => {
    dragging.value = false;
    resizing.value = false;
};

const onResizeDown = (e: PointerEvent) => {
    resizing.value = true;
    const target = e.target as Element | null;
    if (target && typeof (target as HTMLElement).setPointerCapture === 'function') {
        (target as HTMLElement).setPointerCapture(e.pointerId);
    }
    start.value = {
        mx: e.clientX,
        my: e.clientY,
        xMm: frame.value.xMm,
        yMm: frame.value.yMm,
        wMm: frame.value.wMm,
        hMm: frame.value.hMm,
    };
    emit('select');
};

onUnmounted(() => {
    dragging.value = false;
    resizing.value = false;
});
</script>

<template>
    <div :class="containerClass" :style="containerStyle" @pointerdown="onDown" @pointermove="onMove" @pointerup="onUp">
        <div
            v-if="frame.kind === 'text'"
            class="h-full w-full p-1.5 text-[11pt] leading-tight text-gray-800"
            :style="{ transform: `rotate(${frame.rotateDeg ?? -24}deg)` }"
        >
            {{ frame.text }}
        </div>
        <div v-if="frame.kind === 'field'" class="h-full w-full p-1.5 text-[11pt] leading-tight text-gray-600">
            {{ frame.path }}
        </div>
        <img v-if="frame.kind === 'image'" class="h-full w-full object-contain opacity-80" :src="frame.src || '/logo.svg'" />
        <div v-if="frame.kind === 'checkbox'" class="flex h-full w-full items-center gap-2 p-1">
            <div class="border" :style="{ width: '5mm', height: '5mm' }"></div>
            <span class="text-xs">{{ frame.label || 'Checkbox' }}</span>
        </div>
        <div
            v-if="frame.kind === 'grid'"
            class="grid h-full w-full"
            :style="{
                gridTemplateColumns: `repeat(${frame.cols || 8}, 1fr)`,
                gap: `${frame.gapMm || 1}mm`,
            }"
        >
            <div v-for="(_, index) in (frame.cols || 8) * (frame.rows || 1)" :key="index" class="border"></div>
        </div>
        <div v-if="frame.kind === 'line'" class="absolute inset-x-0 bottom-[1.5mm] border-b"></div>
        <div v-if="frame.kind === 'section'" class="h-full w-full p-1.5">
            <div class="mb-1 text-[10pt] font-semibold text-red-600">{{ frame.title || 'Section' }}</div>
        </div>
        <div
            class="bg-primary-600 absolute -bottom-1.5 -right-1.5 h-3 w-3 cursor-nwse-resize rounded-sm opacity-0 group-hover:opacity-100"
            @pointerdown="onResizeDown"
        ></div>
    </div>
</template>

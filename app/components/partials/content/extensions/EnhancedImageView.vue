<script setup lang="ts">
import { getAlignStyle } from '~/composables/tiptap/extensions/EnhancedImage';

const props = defineProps<{
    src: string;
    alt?: string | null;
    title?: string | null;
    width?: number | null;
    height?: number | null;
    align?: string | null;
}>();

const open = ref(false);

const zoom = ref(1);
const fitZoom = ref(1);

const MIN_ZOOM = computed(() => fitZoom.value);
const MAX_ZOOM = 5;

const rotation = ref(0);

const offset = reactive({ x: 0, y: 0 });

const dragging = ref(false);
const hasDragged = ref(false);

const dragStart = reactive({ x: 0, y: 0 });
const offsetStart = reactive({ x: 0, y: 0 });

const imageStyle = computed(() => ({
    transform: `translate(${offset.x}px, ${offset.y}px) rotate(${rotation.value}deg) scale(${zoom.value})`,
    transformOrigin: 'center center',
    cursor: dragging.value ? 'grabbing' : 'grab',
}));

const normalizedRotation = computed(() => {
    return ((rotation.value % 360) + 360) % 360;
});

const isSideways = computed(() => {
    return normalizedRotation.value === 90 || normalizedRotation.value === 270;
});

function onPointerDown(event: PointerEvent) {
    if (event.button !== 0) return;

    dragging.value = true;
    hasDragged.value = false;

    dragStart.x = event.clientX;
    dragStart.y = event.clientY;

    offsetStart.x = offset.x;
    offsetStart.y = offset.y;
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
}

function onPointerMove(event: PointerEvent) {
    if (!dragging.value) return;

    const dx = event.clientX - dragStart.x;
    const dy = event.clientY - dragStart.y;

    if (Math.abs(dx) > 2 || Math.abs(dy) > 2) {
        hasDragged.value = true;
    }

    offset.x = offsetStart.x + dx;
    offset.y = offsetStart.y + dy;
}

function onPointerUp(event: PointerEvent) {
    if (!dragging.value) return;

    dragging.value = false;

    const target = event.currentTarget as HTMLElement;
    if (target.hasPointerCapture(event.pointerId)) {
        target.releasePointerCapture(event.pointerId);
    }
}

function zoomAtPoint(clientX: number, clientY: number, nextZoom: number) {
    const container = previewContainerRef.value;
    if (!container) return;

    nextZoom = Math.min(Math.max(nextZoom, MIN_ZOOM.value), MAX_ZOOM);

    const rect = container.getBoundingClientRect();

    const pointX = clientX - rect.left - rect.width / 2;
    const pointY = clientY - rect.top - rect.height / 2;

    const zoomRatio = nextZoom / zoom.value;

    offset.x = pointX - (pointX - offset.x) * zoomRatio;
    offset.y = pointY - (pointY - offset.y) * zoomRatio;

    zoom.value = nextZoom;

    if (zoom.value <= 1) {
        offset.x = 0;
        offset.y = 0;
    }
}

function zoomInAtPointer(event: MouseEvent) {
    if (hasDragged.value) return;

    zoomAtPoint(event.clientX, event.clientY, zoom.value + 0.25);
}

function zoomOutAtPointer(event: MouseEvent) {
    zoomAtPoint(event.clientX, event.clientY, zoom.value - 0.25);
}

function zoomIn() {
    zoom.value = Math.min(zoom.value + 0.25, MAX_ZOOM);
}

function zoomOut() {
    zoom.value = Math.max(zoom.value - 0.25, MIN_ZOOM.value);

    if (zoom.value <= MIN_ZOOM.value) {
        zoom.value = MIN_ZOOM.value;
        offset.x = 0;
        offset.y = 0;
    }
}

const previewContainerRef = ref<HTMLElement | null>(null);
const previewImageRef = ref<HTMLImageElement | null>(null);

function zoomToFit() {
    const container = previewContainerRef.value;
    const image = previewImageRef.value;

    if (!container || !image) return;

    const rect = container.getBoundingClientRect();

    const naturalWidth = image.naturalWidth;
    const naturalHeight = image.naturalHeight;

    if (!rect.width || !rect.height || !naturalWidth || !naturalHeight) {
        return;
    }

    const sideways = isSideways.value;

    const effectiveWidth = sideways ? naturalHeight : naturalWidth;
    const effectiveHeight = sideways ? naturalWidth : naturalHeight;

    fitZoom.value = Math.min(rect.width / effectiveWidth, rect.height / effectiveHeight, 1);

    zoom.value = fitZoom.value;
    offset.x = 0;
    offset.y = 0;
}

async function rotateLeft() {
    rotation.value -= 90;
    await nextTick();
    zoomToFit();
}

async function rotateRight() {
    rotation.value += 90;
    await nextTick();
    zoomToFit();
}

const alignmentStyle = computed(() => getAlignStyle(props.align));
</script>

<template>
    <UModal
        v-model:open="open"
        :title="title || alt || $t('common.image')"
        fullscreen
        :ui="{ body: 'p-0 sm:p-0 flex flex-col min-h-0 h-[calc(100vh-8rem)]', footer: 'shrink-0' }"
    >
        <img
            class="h-auto max-w-full cursor-zoom-in select-none"
            :src="src"
            :alt="alt || ''"
            :title="title || undefined"
            :width="width || undefined"
            :height="height || undefined"
            :style="[alignmentStyle]"
            @click="open = true"
        />

        <template #body>
            <div
                ref="previewContainerRef"
                class="flex min-h-0 flex-1 items-center justify-center overflow-hidden"
                @pointerdown="onPointerDown"
                @pointermove="onPointerMove"
                @pointerup="onPointerUp"
                @pointercancel="onPointerUp"
                @contextmenu.prevent="zoomOutAtPointer"
            >
                <img
                    ref="previewImageRef"
                    class="block select-none"
                    :class="{ 'transition-transform': !dragging }"
                    :src="src"
                    :alt="alt || ''"
                    :style="imageStyle"
                    draggable="false"
                    @load="zoomToFit"
                    @dragstart.prevent
                    @click.stop="zoomInAtPointer"
                />
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="w-full">
                <UTooltip :text="$t('common.rotate_left')">
                    <UButton class="flex-1" block icon="i-mdi-rotate-left" variant="soft" @click="rotateLeft" />
                </UTooltip>

                <UTooltip :text="$t('common.zoom_out')">
                    <UButton class="flex-1" block icon="i-mdi-zoom-out" variant="soft" @click="zoomOut" />
                </UTooltip>

                <UTooltip :text="$t('common.zoom_to_fit')">
                    <UButton class="flex-1" block icon="i-mdi-fit-to-page-outline" variant="soft" @click="zoomToFit" />
                </UTooltip>

                <UTooltip :text="$t('common.zoom_in')">
                    <UButton class="flex-1" block icon="i-mdi-zoom-in" variant="soft" @click="zoomIn" />
                </UTooltip>

                <UTooltip :text="$t('common.rotate_right')">
                    <UButton class="flex-1" block icon="i-mdi-rotate-right" variant="soft" @click="rotateRight" />
                </UTooltip>
            </UFieldGroup>
        </template>
    </UModal>
</template>

<script setup lang="ts">
import { type NodeViewProps, NodeViewWrapper } from '@tiptap/vue-3';
import { computed, ref, watch } from 'vue';
import AlignmentBar from '~/components/partials/editor/AlignmentBar.vue';

const props = defineProps<NodeViewProps>();

const imgRef = ref<HTMLImageElement | null>(null);
const showAlignmentBar = ref(false);
const isResizing = ref(false);
const startX = ref(0);
const startWidth = ref(0);
const nodeViewRoot = ref<HTMLElement | null>(null);

const currentAlignment = computed<'left' | 'center' | 'right'>(() => {
    const style = (props.node.attrs.style as string) || '';
    if (style.includes('margin: 0 auto 0 0')) return 'left';
    if (style.includes('margin: 0 auto;')) return 'center';
    if (style.includes('margin: 0 0 0 auto')) return 'right';
    return 'left';
});

function setAlignment(align: 'left' | 'center' | 'right') {
    let margin = '';
    if (align === 'left') margin = 'margin: 0 auto 0 0;';
    else if (align === 'center') margin = 'margin: 0 auto;';
    else if (align === 'right') margin = 'margin: 0 0 0 auto;';
    const style = ((props.node.attrs.style as string) || '').replace(/margin:[^;]+;/, '') + ' ' + margin;
    props.updateAttributes({ style });
}

function onImageClick() {
    if (props.editor.isEditable) showAlignmentBar.value = true;
}

function onBlur() {
    showAlignmentBar.value = false;
}

function onResizeMouseDown(e: MouseEvent, corner: number) {
    if (!props.editor.isEditable) return;
    e.preventDefault();
    isResizing.value = true;
    startX.value = e.clientX;
    startWidth.value = imgRef.value?.width || 0;
    document.addEventListener('mousemove', (ev) => onResizeMouseMove(ev, corner));
    document.addEventListener('mouseup', onResizeMouseUp, { once: true });
}

function onResizeMouseMove(e: MouseEvent, corner: number) {
    if (!isResizing.value || !imgRef.value) return;
    const deltaX = corner % 2 === 0 ? -(e.clientX - startX.value) : e.clientX - startX.value;
    const newWidth = Math.max(32, startWidth.value + deltaX);
    imgRef.value.style.width = newWidth + 'px';
    props.updateAttributes({ style: updateStyleWidth(props.node.attrs.style as string, newWidth) });
}

function onResizeMouseUp() {
    isResizing.value = false;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    document.removeEventListener('mousemove', onResizeMouseMove as any);
    // After resizing, if node is not selected, hide the bar
    if (!props.selected) {
        showAlignmentBar.value = false;
    }
}

function onResizeTouchStart(e: TouchEvent, corner: number) {
    if (!props.editor.isEditable) return;
    if (e.cancelable) e.preventDefault();
    isResizing.value = true;
    startX.value = e.touches[0]?.clientX || 0;
    startWidth.value = imgRef.value?.width || 0;
    document.addEventListener('touchmove', (ev) => onResizeTouchMove(ev, corner));
    document.addEventListener('touchend', onResizeTouchEnd, { once: true });
}

function onResizeTouchMove(e: TouchEvent, corner: number) {
    if (!isResizing.value || !imgRef.value || !e.touches[0]) return;
    const deltaX = corner % 2 === 0 ? -(e.touches[0].clientX - startX.value) : e.touches[0].clientX - startX.value;
    const newWidth = Math.max(32, startWidth.value + deltaX);
    imgRef.value.style.width = newWidth + 'px';
    props.updateAttributes({ style: updateStyleWidth(props.node.attrs.style as string, newWidth) });
}

function onResizeTouchEnd() {
    isResizing.value = false;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    document.removeEventListener('touchmove', onResizeTouchMove as any);
    if (!props.selected) {
        showAlignmentBar.value = false;
    }
}

function updateStyleWidth(style: string, width: number) {
    // Remove any previous width and add new
    const s = (style || '').replace(/width:[^;]+;/, '');
    return `width: ${width}px;${s}`;
}

const dotPositions = [
    { top: '-8px', left: '-8px', cursor: 'nwse-resize' },
    { top: '-8px', right: '-8px', cursor: 'nesw-resize' },
    { bottom: '-8px', left: '-8px', cursor: 'nesw-resize' },
    { bottom: '-8px', right: '-8px', cursor: 'nwse-resize' },
];

// Watch the selected prop from Tiptap
watch(
    () => props.selected,
    (selected) => {
        if (!selected && !isResizing.value) {
            showAlignmentBar.value = false;
        }
    },
    { immediate: true },
);
</script>

<template>
    <NodeViewWrapper ref="nodeViewRoot" class="enhanced-image-nodeview" @blur="onBlur">
        <div class="relative w-full">
            <div class="relative" :style="node.attrs.style || ''">
                <img
                    ref="imgRef"
                    v-bind="node.attrs"
                    style="border-radius: 4px"
                    class="h-auto max-w-full cursor-pointer select-none"
                    draggable="true"
                    @click="onImageClick"
                />
                <div
                    v-if="(showAlignmentBar || isResizing) && (selected || isResizing)"
                    class="pointer-events-none absolute -top-4 left-1/2 z-20 w-full -translate-x-1/2 transform"
                >
                    <div class="pointer-events-auto flex justify-center">
                        <AlignmentBar :model-value="currentAlignment" @update:model-value="setAlignment" />
                    </div>
                </div>
                <template v-if="(showAlignmentBar || isResizing) && props.editor.isEditable && (selected || isResizing)">
                    <div
                        v-for="(pos, i) in dotPositions"
                        :key="i"
                        class="resize-dot absolute z-10 h-4 w-4 rounded-full border border-neutral-500 bg-white"
                        :style="{
                            ...pos,
                            borderWidth: '1.5px', // closest to Tailwind border-2
                        }"
                        @mousedown="(e) => onResizeMouseDown(e, i)"
                        @touchstart="(e) => onResizeTouchStart(e, i)"
                    ></div>
                </template>
            </div>
        </div>
    </NodeViewWrapper>
</template>

<style scoped>
.enhanced-image-nodeview img {
    max-width: 100%;
    height: auto;
    border-radius: 4px;
}
.resize-dot {
    box-sizing: border-box;
    transition: background 0.2s;
}
.resize-dot:hover {
    background: #e5e7eb;
}
</style>

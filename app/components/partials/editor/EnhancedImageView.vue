<script setup lang="ts">
import { type NodeViewProps, NodeViewWrapper } from '@tiptap/vue-3';
import { computed, ref, watch } from 'vue';
import AlignmentBar from '~/components/partials/editor/AlignmentBar.vue';
import { getAlignStyle, removeStyleProperties, type ImageAlign } from '~/composables/tiptap/extensions/EnhancedImage';

const props = defineProps<NodeViewProps>();

const imgRef = ref<HTMLImageElement | null>(null);
const showAlignmentBar = ref<boolean>(false);
const isResizing = ref<boolean>(false);
const startX = ref(0);
const startWidth = ref(0);

const currentAlignment = computed<ImageAlign>(() => {
    const align = props.node.attrs.align;
    return align === 'center' || align === 'right' ? align : 'left';
});

function setAlignment(align: ImageAlign) {
    props.updateAttributes({
        align: align,
        style: removeStyleProperties(props.node.attrs.style ?? '', ['margin']),
    });
}

function onImageClick() {
    if (props.editor.isEditable) showAlignmentBar.value = true;
}

function onBlur() {
    showAlignmentBar.value = false;
}

let activeMouseMoveHandler: ((e: MouseEvent) => void) | null = null;

function onResizeMouseDown(e: MouseEvent, corner: number) {
    if (!props.editor.isEditable) return;

    e.preventDefault();

    isResizing.value = true;
    startX.value = e.clientX;
    startWidth.value = imgRef.value?.width || 0;

    activeMouseMoveHandler = (ev) => onResizeMouseMove(ev, corner);

    document.addEventListener('mousemove', activeMouseMoveHandler);
    document.addEventListener('mouseup', onResizeMouseUp, { once: true });
}

function onResizeMouseMove(e: MouseEvent, corner: number) {
    if (!isResizing.value || !imgRef.value) return;

    const deltaX = corner % 2 === 0 ? -(e.clientX - startX.value) : e.clientX - startX.value;
    const newWidth = Math.max(32, startWidth.value + deltaX);

    props.updateAttributes({
        width: newWidth,
        height: null,
    });
}

function onResizeMouseUp() {
    isResizing.value = false;

    if (activeMouseMoveHandler) {
        document.removeEventListener('mousemove', activeMouseMoveHandler);
        activeMouseMoveHandler = null;
    }

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

    props.updateAttributes({
        width: newWidth,
        height: null,
    });
}

function onResizeTouchEnd() {
    isResizing.value = false;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    document.removeEventListener('touchmove', onResizeTouchMove as any);
    if (!props.selected) {
        showAlignmentBar.value = false;
    }
}

const dotPositions = [
    { top: '-8px', left: '-8px', cursor: 'nwse-resize' },
    { top: '-8px', right: '-8px', cursor: 'nesw-resize' },
    { bottom: '-8px', left: '-8px', cursor: 'nesw-resize' },
    { bottom: '-8px', right: '-8px', cursor: 'nwse-resize' },
];

const alignmentStyle = computed(() => getAlignStyle(props.node.attrs.align));
const style = computed(() => removeStyleProperties(props.node.attrs.style ?? '', ['width', 'height', 'margin']));

const imageWidth = computed(() => {
    const width = Number(props.node.attrs.width);
    return Number.isFinite(width) && width > 0 ? width : undefined;
});

const imageHeight = computed(() => {
    const height = Number(props.node.attrs.height);
    return Number.isFinite(height) && height > 0 ? height : undefined;
});

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
    <NodeViewWrapper class="enhanced-image-nodeview relative my-2 w-full" @blur="onBlur">
        <div class="relative w-full">
            <div class="relative w-fit" :style="alignmentStyle">
                <div class="relative inline-block">
                    <img
                        ref="imgRef"
                        class="h-auto max-w-full cursor-pointer select-none"
                        :class="[selected || isResizing ? 'border border-primary-500' : '']"
                        :style="style"
                        :draggable="true"
                        :src="cleanupImageURL(node.attrs.src)"
                        :alt="node.attrs.alt || ''"
                        :title="node.attrs.title || undefined"
                        :width="imageWidth"
                        :height="imageHeight"
                        :data-file-id="node.attrs.fileId || undefined"
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
                            class="resize-dot absolute z-10 h-4 w-4 rounded-full bg-elevated ring ring-accented ring-inset hover:bg-accented/75"
                            :style="{
                                ...pos,
                            }"
                            @mousedown="(e) => onResizeMouseDown(e, i)"
                            @touchstart="(e) => onResizeTouchStart(e, i)"
                        />
                    </template>
                </div>
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
</style>

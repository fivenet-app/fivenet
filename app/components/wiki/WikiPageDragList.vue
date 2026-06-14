<script lang="ts" setup>
import { useDraggable, type UseDraggableOptions } from 'vue-draggable-plus';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

const props = withDefaults(
    defineProps<{
        disabled?: boolean;
        handle?: string;
        draggable?: string;
        animation?: number;
        onMove?: UseDraggableOptions<PageShort>['onMove'];
        onEnd?: UseDraggableOptions<PageShort>['onEnd'];
    }>(),
    {
        disabled: false,
        handle: '.wiki-page-drag-handle',
        draggable: '> .wiki-page-list-item',
        animation: 150,
        onMove: undefined,
        onEnd: undefined,
    },
);

const model = defineModel<PageShort[]>({ required: true });

const attrs = useAttrs();
const listRef = useTemplateRef('listRef');

const draggable = useDraggable(listRef, model, {
    immediate: true,
    animation: props.animation,
    handle: props.handle,
    draggable: props.draggable,
    onMove: props.onMove,
    onEnd: props.onEnd,
});

watch(
    () => props.disabled,
    (disabled) => {
        if (disabled) {
            draggable.pause();
        } else {
            draggable.resume();
        }
    },
    { immediate: true },
);

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <div ref="listRef" v-bind="attrs">
        <slot />
    </div>
</template>

<script setup lang="ts">
import type { WidgetFrame } from './types';
import FieldRenderer from './widgets/FieldRenderer.vue';

const props = defineProps<{
    node: WidgetFrame; // <- narrowed type
    mode: 'layout' | 'fill' | 'preview';
    data: Record<string, unknown>; // form data
    setData: (path: string, v: unknown) => void;
}>();

const value = computed(() => {
    const path = props.node.binding?.path;
    return path ? (props.data[path] ?? '') : '';
});
</script>

<template>
    <div
        class="absolute"
        :style="{ left: node.xMm + 'mm', top: node.yMm + 'mm', width: node.wMm + 'mm', height: node.hMm + 'mm' }"
    >
        <template v-if="mode === 'fill'">
            <FieldRenderer
                :widget="node.widget"
                :value="value"
                :binding="node.binding"
                :rows="node.rows"
                :box-char-count="node.boxCharCount"
                :max-chars="node.maxChars"
                @update="(v) => node.binding?.path && setData(node.binding.path, v)"
            />
        </template>

        <template v-else>
            <!-- layout/preview: draw text as it will print (fit/truncate/shrink) -->
            <RenderFittedText
                v-if="node.binding"
                :text="String(value ?? '')"
                :fit="node.binding.fit"
                :min-font-size="node.binding.minFontSize ?? 9"
            />
            <!-- (optional) ghost box when not bound -->
            <div v-else class="h-full w-full border border-dashed opacity-40"></div>
        </template>
    </div>
</template>

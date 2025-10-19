<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
const props = defineProps<{ api?: any | null }>();
</script>

<template>
    <UCard>
        <template #header>Shapes</template>
        <div class="grid grid-cols-3 gap-2">
            <UButton
                size="sm"
                @click="
                    props.api?.addSvgElementFromJson?.({
                        element: 'rect',
                        attr: { x: 50, y: 50, width: 120, height: 80, fill: '#4ade80' },
                    })
                "
                >Rect</UButton
            >
            <UButton
                size="sm"
                @click="
                    props.api?.addSvgElementFromJson?.({
                        element: 'ellipse',
                        attr: { cx: 120, cy: 120, rx: 60, ry: 40, fill: '#60a5fa' },
                    })
                "
                >Ellipse</UButton
            >
            <UButton
                size="sm"
                @click="
                    props.api?.addSvgElementFromJson?.({
                        element: 'line',
                        attr: { x1: 30, y1: 30, x2: 200, y2: 200, stroke: '#000', 'stroke-width': 2 },
                    })
                "
                >Line</UButton
            >
            <UButton
                size="sm"
                @click="
                    props.api?.addSvgElementFromJson?.({
                        element: 'text',
                        attr: { x: 80, y: 80, fill: '#000', 'font-size': 24, text: 'Hello' },
                    })
                "
                >Text</UButton
            >
        </div>
    </UCard>
</template>

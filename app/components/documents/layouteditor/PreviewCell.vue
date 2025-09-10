<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import type { WidgetFrame } from './types'; // adjust path

const props = defineProps<{
    node: WidgetFrame; // the widget frame to render
    data: Record<string, unknown>; // filled form values
}>();

/** 1) Resolve the bound value and apply simple formatting */
const rawValue = computed(() => {
    const path = props.node.binding?.path;
    const value = (path ? props.data[path] : '') as string | number | null | undefined;
    return value ?? '';
});

const formatted = computed(() => {
    let v = String(rawValue.value ?? '');
    const max = props.node.maxChars;
    if (typeof max === 'number' && max > 0) v = v.slice(0, max);

    const fmt = props.node.binding?.formatter;
    if (!fmt) return v;

    if (fmt === 'upper') return v.toUpperCase();
    if (fmt === 'lower') return v.toLowerCase();
    if (fmt.startsWith('date:')) {
        // naive example: date:DD.MM.YYYY
        try {
            const d = new Date(v);
            const [_, pat] = fmt.split(':');
            if (Number.isFinite(d.getTime())) {
                const dd = `${d.getDate()}`.padStart(2, '0');
                const mm = `${d.getMonth() + 1}`.padStart(2, '0');
                const yyyy = d.getFullYear();
                return pat === 'DD.MM.YYYY' ? `${dd}.${mm}.${yyyy}` : d.toISOString();
            }
        } catch {
            // ignore parse errors
        }
        return v;
    }
    return v;
});

/** 2) Shared box metrics (mm → CSS) */
const styleBox = computed(() => ({
    left: props.node.xMm + 'mm',
    top: props.node.yMm + 'mm',
    width: props.node.wMm + 'mm',
    height: props.node.hMm + 'mm',
}));

/** 3) Fit strategy for text/textarea */
const host = ref<HTMLElement | null>(null);
const fontSizePx = ref<number | null>(null);
const ro = new ResizeObserver(() => recalcShrink());

onMounted(() => {
    if (host.value) ro.observe(host.value);
});
onBeforeUnmount(() => ro.disconnect());

watch([formatted, () => props.node.binding?.fit, () => props.node.binding?.minFontSize], () => {
    fontSizePx.value = null;
    // allow DOM to update, then recompute:
    queueMicrotask(recalcShrink);
});

function measureTextWidth(text: string, font: string): number {
    const c = document.createElement('canvas');
    const ctx = c.getContext('2d')!;
    ctx.font = font;
    return ctx.measureText(text).width;
}

function recalcShrink() {
    if (!host.value) return;
    if (props.node.binding?.fit !== 'shrink') return;

    // Only shrink single-line cases; for multi-line, you can enhance later.
    const el = host.value;
    const text = formatted.value || '';
    const cs = getComputedStyle(el);
    const targetWidth = el.clientWidth - (parseFloat(cs.paddingLeft) + parseFloat(cs.paddingRight) + 1);
    if (targetWidth <= 0) return;

    // binary search font-size
    const minPx = props.node.binding?.minFontSize ?? 9;
    let low = minPx,
        high = parseFloat(cs.fontSize) || 14,
        best = low;

    // ensure we try at least current
    if (high < low) high = low;

    // construct a font string while mutating size
    const baseFont = `${cs.fontStyle} ${cs.fontVariant} ${cs.fontWeight} __SIZE__/${cs.lineHeight} ${cs.fontFamily}`;

    for (let i = 0; i < 16; i++) {
        const mid = Math.floor((low + high) / 2);
        const font = baseFont.replace('__SIZE__', `${mid}px`);
        const w = measureTextWidth(text, font);
        if (w <= targetWidth) {
            best = mid;
            // try larger
            low = mid + 1;
        } else {
            // too big
            high = mid - 1;
        }
    }
    fontSizePx.value = best;
}

/** 4) Boxed-text helpers */
const boxCount = computed(() => {
    const n = props.node.boxCharCount ?? props.node.maxChars ?? 10;
    return Math.max(1, Math.min(128, n));
});
const boxedChars = computed(() => {
    const s = (formatted.value ?? '').toString();
    const arr = Array.from({ length: boxCount.value }, (_, i) => s[i] ?? '');
    return arr;
});
</script>

<template>
    <div class="absolute" :style="styleBox">
        <!-- TEXT / TEXTAREA PREVIEW -->
        <div
            v-if="node.widget === 'text' || node.widget === 'textarea'"
            ref="host"
            class="h-full w-full"
            :style="{
                overflow: node.binding?.fit === 'truncate' ? 'hidden' : 'auto',
                whiteSpace: node.widget === 'text' ? (node.binding?.fit === 'wrap' ? 'pre-wrap' : 'nowrap') : 'pre-wrap',
                textOverflow: node.binding?.fit === 'truncate' ? 'ellipsis' : 'clip',
                fontSize: fontSizePx !== null ? fontSizePx + 'px' : undefined,
                lineHeight: '1.2',
                padding: '1mm',
            }"
        >
            {{ formatted }}
        </div>

        <!-- BOXED TEXT PREVIEW -->
        <div v-else-if="node.widget === 'boxed-text'" class="flex h-full w-full items-center" :style="{ padding: '1mm' }">
            <div
                class="grid w-full"
                :style="{
                    gridTemplateColumns: `repeat(${boxCount}, 1fr)`,
                    gap: (node.letterSpacing ?? 0) + 'px',
                }"
            >
                <div
                    v-for="(ch, i) in boxedChars"
                    :key="i"
                    class="flex items-center justify-center border"
                    :style="{ aspectRatio: '7/9' }"
                >
                    <span class="leading-none select-none">{{ ch }}</span>
                </div>
            </div>
        </div>

        <!-- FALLBACK (shouldn't happen if widget-only) -->
        <div v-else class="h-full w-full border border-dashed opacity-40"></div>
    </div>
</template>

<style scoped>
.border {
    border-color: #999;
    border-width: 1px;
}
</style>

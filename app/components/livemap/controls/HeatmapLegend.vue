<script setup lang="ts">
const props = withDefaults(
    defineProps<{
        // Whether to show the legend
        show?: boolean;
        max?: number;
        // Optional color ramp; keys 0‒1 → color strings
        gradient?: Record<number, string>;
    }>(),
    {
        show: true,
        max: 0,
        gradient: undefined,
    },
);

// Default ramp resembles the example used for L.heatLayer
const gradient = computed<Record<number, string>>(
    () => props.gradient ?? { 0.2: 'blue', 0.4: 'lime', 0.6: 'orange', 0.8: 'red' },
);

// CSS style for the colored bar
const barStyle = computed(() => {
    const stops = Object.entries(gradient.value)
        .sort((a, b) => Number(a[0]) - Number(b[0]))
        .map(([stop, color]) => `${color} ${Number(stop) * 100}%`)
        .join(',');
    return {
        background: `linear-gradient(to right, ${stops})`,
    };
});
</script>

<template>
    <LControl position="bottomright">
        <div v-if="show" class="space-y-1 rounded-md bg-default p-1 text-xs font-medium shadow-md">
            <!-- Gradient bar -->
            <div class="h-2 w-32 rounded-sm" :style="barStyle" />
            <!-- Captions -->
            <div class="flex justify-between text-toned">
                <span>{{ $t('common.min') }}</span>
                <span v-if="max > 0" class="text-default">{{ max }}</span>
                <span>{{ $t('common.max') }}</span>
            </div>
        </div>
    </LControl>
</template>

<script setup lang="ts">
import type { PointTuple } from 'leaflet';
import { computed } from 'vue';

const props = withDefaults(
    defineProps<{
        latLng: [number, number];
        radius?: number;
        zoom?: number;

        fillOpacity?: number;
        zIndexOffset?: number;

        baseZoom?: number;
        minMarkerScale?: number;
        maxMarkerScale?: number;
        zoomShrinkPower?: number;

        color?: string;
        pulseDuration?: number;
        pulseDelay?: number;
        ringSize?: number;
        ringGap?: number;
        ringWidth?: number;
    }>(),
    {
        radius: 100,
        zoom: 4,

        fillOpacity: 0.38,
        zIndexOffset: 1000,

        baseZoom: 5,
        minMarkerScale: 0.01,
        maxMarkerScale: 1,
        zoomShrinkPower: 1.7,

        color: '#e83e8c',
        pulseDuration: 2.0,
        pulseDelay: 0.45,
        ringSize: 64,
        ringGap: 16,
        ringWidth: 4,
    },
);

const markerSize = computed(() => {
    return props.ringSize + props.ringGap * 2 + props.ringWidth * 2;
});

const iconSize = computed<PointTuple>(() => [markerSize.value, markerSize.value]);
const iconAnchor = computed<PointTuple>(() => [markerSize.value / 2, markerSize.value / 2]);

const markerScale = computed(() => {
    const zoomDelta = props.zoom - props.baseZoom;

    const scale = zoomDelta >= 0 ? 1 : 2 ** (zoomDelta * props.zoomShrinkPower);

    return Math.min(props.maxMarkerScale, Math.max(props.minMarkerScale, scale));
});

const cssVars = computed(() => ({
    '--zone-color': props.color,
    '--pulse-duration': `${props.pulseDuration}s`,
    '--pulse-delay': `${props.pulseDelay}s`,
    '--ring-size': `${props.ringSize}px`,
    '--ring-gap': `${props.ringGap}px`,
    '--ring-width': `${props.ringWidth}px`,
    '--marker-scale': String(markerScale.value),
}));
</script>

<template>
    <!-- Center pulsing marker -->
    <LMarker :lat-lng="latLng" :z-index-offset="100" :options="{ interactive: false }">
        <LIcon class-name="pointer-events-none!" :icon-size="iconSize" :icon-anchor="iconAnchor">
            <div class="zone-radar-marker" :style="cssVars">
                <div class="zone-radar-scale">
                    <span class="zone-radar-ring ring-1" />
                    <span class="zone-radar-ring ring-2" />
                    <span class="zone-radar-ring ring-3" />
                </div>
            </div>
        </LIcon>
    </LMarker>
</template>

<style scoped>
:deep(.leaflet-div-icon),
:deep(.leaflet-marker-icon) {
    background: transparent !important;
    border: none !important;
    box-shadow: none !important;
}

.zone-radar-marker {
    position: relative;
    width: 100%;
    height: 100%;
    pointer-events: none;
}

.zone-radar-scale {
    position: absolute;
    inset: 0;
    transform: scale(var(--marker-scale));
    transform-origin: center;
}

.zone-radar-ring {
    position: absolute;
    left: 50%;
    top: 50%;

    border: var(--ring-width) solid var(--zone-color);
    border-radius: 999px;
    box-sizing: border-box;

    background: transparent;
    outline: none;
    box-shadow: none;

    transform: translate(-50%, -50%);
    opacity: 0.18;

    animation: zone-radar-ring-pulse var(--pulse-duration) infinite linear;
}

.ring-1 {
    width: var(--ring-size);
    height: var(--ring-size);
    animation-delay: 0s;
}

.ring-2 {
    width: calc(var(--ring-size) + var(--ring-gap));
    height: calc(var(--ring-size) + var(--ring-gap));
    animation-delay: var(--pulse-delay);
}

.ring-3 {
    width: calc(var(--ring-size) + var(--ring-gap) * 2);
    height: calc(var(--ring-size) + var(--ring-gap) * 2);
    animation-delay: calc(var(--pulse-delay) * 2);
}

@keyframes zone-radar-ring-pulse {
    0% {
        opacity: 0.18;
    }

    6% {
        opacity: 1;
    }

    12% {
        opacity: 0.68;
    }

    18% {
        opacity: 0.24;
    }

    100% {
        opacity: 0.18;
    }
}
</style>

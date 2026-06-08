<script lang="ts" setup>
import type { LatLngBoundsExpression } from 'leaflet';

defineProps<{
    tileLayer: string;
    bounds: LatLngBoundsExpression;
}>();
</script>

<template>
    <LImageOverlay :url="`/images/livemap/overlays/cayo-perico/${tileLayer}.webp`" :bounds="bounds" />
</template>

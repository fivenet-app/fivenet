<script lang="ts" setup>
import { LCircleMarker, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import L from 'leaflet';
import { Marker } from '~~/gen/ts/resources/livemap/livemap';

const props = withDefaults(
    defineProps<{
        marker: Marker;
        size?: number;
    }>(),
    {
        size: 20,
    },
);

defineEmits<{
    (e: 'selected'): void;
}>();

const iconAnchor: L.PointExpression = [props.size / 2, props.size];
const popupAnchor: L.PointExpression = [0, (props.size / 2) * -1];
const icon = new L.DivIcon({
    html: `<div>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -0.8 16 17.6" fill="${
            props.marker.info?.color ? '#' + props.marker.info?.color : 'currentColor'
        }" class="w-full h-full">
                <path d="M8 16s6-5.686 6-10A6 6 0 0 0 2 6c0 4.314 6 10 6 10zm0-7a3 3 0 1 1 0-6 3 3 0 0 1 0 6z"/>
            </svg>
        </div>`,
    iconSize: [props.size, props.size],
    iconAnchor,
    popupAnchor,
}) as L.Icon;
</script>

<template>
    <LCircleMarker
        v-if="marker.data?.data.oneofKind === 'circle'"
        :key="marker.info!.id?.toString()"
        :latLng="[marker.info!.y, marker.info!.x]"
        :radius="marker.data?.data.circle.radius"
    >
    </LCircleMarker>

    <LMarker
        v-else
        :latLng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        :icon="icon"
        @click="$emit('selected')"
    >
        <LPopup :options="{ closeButton: true }">
            {{ marker.info?.name }}
        </LPopup>
    </LMarker>
</template>

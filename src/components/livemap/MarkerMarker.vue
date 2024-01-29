<script lang="ts" setup>
import { LCircle, LIcon, LMarker } from '@vue-leaflet/vue-leaflet';
import { type PointExpression } from 'leaflet';
import { MapMarkerQuestionIcon } from 'mdi-vue3';
import { Marker } from '~~/gen/ts/resources/livemap/livemap';
import { markerIcons } from '~/components/livemap/helpers';
import MarkerMarkerPopup from '~/components/livemap/MarkerMarkerPopup.vue';

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

const iconAnchor: PointExpression = [props.size / 2, props.size];
const popupAnchor: PointExpression = [0, (props.size / 2) * -1];
</script>

<template>
    <LCircle
        v-if="marker.data?.data.oneofKind === 'circle'"
        :key="marker.info!.id"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :radius="marker.data?.data.circle.radius / 0.6931471805599453"
        :color="marker.info?.color ? '#' + marker.info?.color : '#ffffff'"
        :fill-color="marker.info?.color ? '#' + marker.info?.color : '#ffffff'"
        :fill-opacity="(marker.data.data.circle.oapcity ?? 15) / 100"
    >
        <MarkerMarkerPopup :marker="marker" />
    </LCircle>

    <LMarker
        v-else-if="marker.data?.data.oneofKind === 'icon'"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        @click="$emit('selected')"
    >
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <component
                :is="
                    markerIcons.find((i) => marker.data?.data.oneofKind === 'icon' && i.name === marker.data?.data.icon.icon) ??
                    MapMarkerQuestionIcon
                "
                class="h-6 w-6"
                :style="{ color: marker.info?.color ? '#' + marker.info?.color : 'currentColor' }"
            />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>

    <LMarker v-else :lat-lng="[marker.info!.y, marker.info!.x]" :name="marker.info!.name" @click="$emit('selected')">
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <MapMarkerQuestionIcon :fill="marker.info?.color ? '#' + marker.info?.color : 'currentColor'" class="h-5 w-5" />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>
</template>

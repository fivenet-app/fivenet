<script lang="ts" setup>
import { LCircle, LIcon, LMarker } from '@vue-leaflet/vue-leaflet';
import type { PointExpression } from 'leaflet';
import MarkerMarkerPopup from '~/components/livemap/MarkerMarkerPopup.vue';
import { availableIcons, fallbackIcon } from '~/components/partials/icons';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/livemap';

const props = withDefaults(
    defineProps<{
        marker: MarkerMarker;
        size?: number;
    }>(),
    {
        size: 20,
    },
);

defineEmits<{
    (e: 'selected'): void;
}>();

const { livemap } = useAppConfig();

const iconAnchor = ref<PointExpression>([props.size / 2, props.size]);
const popupAnchor = ref<PointExpression>([0, (props.size / 2) * -1]);
</script>

<template>
    <LCircle
        v-if="marker.data?.data.oneofKind === 'circle'"
        :key="marker.info!.id"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :radius="marker.data?.data.circle.radius / 0.6931471805599453"
        :color="marker.info?.color ?? livemap.markerMarkers.fallbackColor"
        :fill-color="marker.info?.color ?? livemap.markerMarkers.fallbackColor"
        :fill-opacity="(marker.data.data.circle.opacity ?? 15) / 100"
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
                    availableIcons.find(
                        (icon) =>
                            marker.data?.data.oneofKind === 'icon' &&
                            icon.name === convertDynamicIconNameToComponent(marker.data?.data.icon.icon),
                    ) ?? fallbackIcon.name
                "
                class="size-full"
                :style="{ color: marker.info?.color ?? 'currentColor' }"
            />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>

    <LMarker v-else :lat-lng="[marker.info!.y, marker.info!.x]" :name="marker.info!.name" @click="$emit('selected')">
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <component :is="fallbackIcon" :fill="marker.info?.color ?? 'currentColor'" />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>
</template>

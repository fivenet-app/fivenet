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
        :key="marker.id"
        :name="marker.name"
        :lat-lng="[marker.y, marker.x]"
        :radius="marker.data?.data.circle.radius / 0.6931471805599453"
        :color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-opacity="(marker.data.data.circle.opacity ?? 15) / 100"
    >
        <MarkerMarkerPopup :marker="marker" />
    </LCircle>

    <LMarker
        v-else-if="marker.data?.data.oneofKind === 'icon'"
        :name="marker.name"
        :lat-lng="[marker.y, marker.x]"
        @click="$emit('selected')"
    >
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <component
                :is="
                    availableIcons.find(
                        (icon) =>
                            marker.data?.data.oneofKind === 'icon' &&
                            icon.name === convertDynamicIconNameToComponent(marker.data?.data.icon.icon),
                    )?.component ?? fallbackIcon.name
                "
                class="size-full"
                :style="{ color: marker.color ?? 'currentColor' }"
            />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>

    <LMarker v-else :name="marker.name" :lat-lng="[marker.y, marker.x]" @click="$emit('selected')">
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <component :is="fallbackIcon" :fill="marker.color ?? 'currentColor'" />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>
</template>

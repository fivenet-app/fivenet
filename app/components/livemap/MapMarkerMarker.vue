<script lang="ts" setup>
import type { DragEndEvent, LatLngTuple, LeafletEvent, PointExpression } from 'leaflet';
import MarkerMarkerPopup from '~/components/livemap/MarkerMarkerPopup.vue';
import MarkerCreateOrUpdateSlideover from '~/components/livemap/MarkerCreateOrUpdateSlideover.vue';
import { resolveIconComponent } from '~/components/partials/icons';
import { useLivemapStore } from '~/stores/livemap';
import { MarkerType, type MarkerMarker } from '~~/gen/ts/resources/livemap/markers/marker_marker';

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

const { can, activeChar, isSuperuser } = useAuth();

const { livemap } = useAppConfig();
const overlay = useOverlay();
const livemapStore = useLivemapStore();
const { suppressMapPreclick } = livemapStore;
const { markerDragEnabled, markersMarkers } = storeToRefs(livemapStore);
const markerCreateOrUpdateSlideover = overlay.create(MarkerCreateOrUpdateSlideover);

const iconAnchor = computed<PointExpression>(() => [props.size / 2, props.size]);
const popupAnchor = computed<PointExpression>(() => [0, (props.size / 2) * -1]);
const dragHandleAnchor: PointExpression = [5, 5];
const dragHandleSize: PointExpression = [10, 10];
const canEditMarkers = computed(() => can('livemap.LivemapService/CreateOrUpdateMarker').value);

function cloneMarkerPlain(marker: MarkerMarker): MarkerMarker {
    return JSON.parse(JSON.stringify(toRaw(marker))) as MarkerMarker;
}

function canMoveMarker(marker: MarkerMarker): boolean {
    if (!markerDragEnabled.value) return false;
    if (!canEditMarkers.value || marker.id <= 0) return false;
    if (isSuperuser.value) return true;
    return !!activeChar.value?.job && marker.job === activeChar.value.job;
}

function hasPersistedPoint(markerId: number, x: number, y: number): boolean {
    const persisted = markersMarkers.value.get(markerId);
    if (!persisted) return false;
    const epsilon = 0.000001;
    return Math.abs(persisted.x - x) <= epsilon && Math.abs(persisted.y - y) <= epsilon;
}

function hasPersistedRectangleEnd(markerId: number, endX: number, endY: number): boolean {
    const persisted = markersMarkers.value.get(markerId);
    if (persisted?.data?.data.oneofKind !== 'rectangle') return false;
    const epsilon = 0.000001;
    return (
        Math.abs(persisted.data.data.rectangle.endX - endX) <= epsilon &&
        Math.abs(persisted.data.data.rectangle.endY - endY) <= epsilon
    );
}

function hasPersistedPolygonPoint(markerId: number, pointIndex: number, x: number, y: number): boolean {
    const persisted = markersMarkers.value.get(markerId);
    if (persisted?.data?.data.oneofKind !== 'polygon') return false;
    const epsilon = 0.000001;

    if (pointIndex === 0) {
        return Math.abs(persisted.x - x) <= epsilon && Math.abs(persisted.y - y) <= epsilon;
    }

    const point = persisted.data.data.polygon.points[pointIndex - 1];
    if (!point) return false;
    return Math.abs(point.x - x) <= epsilon && Math.abs(point.y - y) <= epsilon;
}

function openDragEditSlideover(nextMarker: MarkerMarker, onClose: (saved: unknown) => void): void {
    markerCreateOrUpdateSlideover.open({
        marker: nextMarker,
        onClose: (saved) => onClose(saved),
    });
}

function onMarkerDragStart(event: LeafletEvent): void {
    event.target.closePopup?.();
    suppressMapPreclick(1000);
}

async function onMarkerDragEnd(event: DragEndEvent, marker: MarkerMarker): Promise<void> {
    if (!canMoveMarker(marker)) return;

    suppressMapPreclick(500);

    const latlng = event.target.getLatLng();
    const previousLatLng: LatLngTuple = [marker.y, marker.x];
    const nextMarker = cloneMarkerPlain(marker);
    nextMarker.x = latlng.lng;
    nextMarker.y = latlng.lat;

    openDragEditSlideover(nextMarker, (saved) => {
        if (saved === true || hasPersistedPoint(marker.id, latlng.lng, latlng.lat)) return;
        event.target.setLatLng(previousLatLng);
    });
}

async function onRectanglePointDragEnd(event: DragEndEvent, marker: MarkerMarker, startPoint: boolean): Promise<void> {
    if (!canMoveMarker(marker) || marker.data?.data.oneofKind !== 'rectangle') return;

    suppressMapPreclick(500);

    const latlng = event.target.getLatLng();
    const previousLatLng: [number, number] = startPoint
        ? [marker.y, marker.x]
        : [marker.data.data.rectangle.endY, marker.data.data.rectangle.endX];

    const nextMarker = cloneMarkerPlain(marker);
    if (nextMarker.data?.data.oneofKind !== 'rectangle') return;

    if (startPoint) {
        nextMarker.x = latlng.lng;
        nextMarker.y = latlng.lat;
    } else {
        nextMarker.data.data.rectangle.endX = latlng.lng;
        nextMarker.data.data.rectangle.endY = latlng.lat;
    }

    openDragEditSlideover(nextMarker, (saved) => {
        if (
            saved === true ||
            (startPoint
                ? hasPersistedPoint(marker.id, latlng.lng, latlng.lat)
                : hasPersistedRectangleEnd(marker.id, latlng.lng, latlng.lat))
        ) {
            return;
        }
        event.target.setLatLng(previousLatLng);
    });
}

async function onPolygonPointDragEnd(event: DragEndEvent, marker: MarkerMarker, pointIndex: number): Promise<void> {
    if (!canMoveMarker(marker) || marker.data?.data.oneofKind !== 'polygon') return;

    suppressMapPreclick(500);

    const latlng = event.target.getLatLng();
    const previousLatLng: [number, number] =
        pointIndex === 0
            ? [marker.y, marker.x]
            : [marker.data.data.polygon.points[pointIndex - 1]!.y, marker.data.data.polygon.points[pointIndex - 1]!.x];

    const nextMarker = cloneMarkerPlain(marker);
    if (nextMarker.data?.data.oneofKind !== 'polygon') return;

    if (pointIndex === 0) {
        nextMarker.x = latlng.lng;
        nextMarker.y = latlng.lat;
    } else {
        nextMarker.data.data.polygon.points[pointIndex - 1] = {
            x: latlng.lng,
            y: latlng.lat,
        };
    }

    openDragEditSlideover(nextMarker, (saved) => {
        if (saved === true || hasPersistedPolygonPoint(marker.id, pointIndex, latlng.lng, latlng.lat)) return;
        event.target.setLatLng(previousLatLng);
    });
}

function getPolygonDragPoints(marker: MarkerMarker): { x: number; y: number }[] {
    if (marker.data?.data.oneofKind !== 'polygon') return [];
    return [{ x: marker.x, y: marker.y }, ...marker.data.data.polygon.points];
}
</script>

<template>
    <LMarker
        v-if="marker.type === MarkerType.DOT"
        :name="marker.name"
        :lat-lng="[marker.y, marker.x]"
        :draggable="canMoveMarker(marker)"
        :options="{ markerMarker: marker }"
        @click="$emit('selected')"
        @dragstart="onMarkerDragStart"
        @dragend="onMarkerDragEnd($event, marker)"
    >
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <div class="size-full" style="background-color: white; border-radius: 50%; border: 1px solid #000"></div>
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>

    <LCircle
        v-if="marker.data?.data.oneofKind === 'circle'"
        :key="marker.id"
        :name="marker.name"
        :lat-lng="[marker.y, marker.x]"
        :radius="marker.data?.data.circle.radius / 0.6931471805599453"
        :color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-opacity="(marker.data.data.circle.opacity ?? 15) / 100"
        :options="{ markerMarker: marker }"
    >
        <MarkerMarkerPopup :marker="marker" />
    </LCircle>
    <LMarker
        v-if="marker.data?.data.oneofKind === 'circle' && canMoveMarker(marker)"
        :lat-lng="[marker.y, marker.x]"
        :draggable="true"
        @dragstart="onMarkerDragStart"
        @dragend="onMarkerDragEnd($event, marker)"
    >
        <LIcon :icon-size="dragHandleSize" :icon-anchor="dragHandleAnchor">
            <div class="size-full rounded-full border border-black bg-white/90"></div>
        </LIcon>
    </LMarker>

    <LRectangle
        v-if="marker.data?.data.oneofKind === 'rectangle'"
        :name="marker.name"
        :bounds="[
            [marker.y, marker.x],
            [marker.data.data.rectangle.endY, marker.data.data.rectangle.endX],
        ]"
        :color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-opacity="(marker.data.data.rectangle.opacity ?? 15) / 100"
        :options="{ markerMarker: marker }"
    >
        <MarkerMarkerPopup :marker="marker" />
    </LRectangle>
    <LMarker
        v-if="marker.data?.data.oneofKind === 'rectangle' && canMoveMarker(marker)"
        :lat-lng="[marker.y, marker.x]"
        :draggable="true"
        @dragstart="onMarkerDragStart"
        @dragend="onRectanglePointDragEnd($event, marker, true)"
    >
        <LIcon :icon-size="dragHandleSize" :icon-anchor="dragHandleAnchor">
            <div class="size-full rounded-sm border border-black bg-white/90"></div>
        </LIcon>
    </LMarker>
    <LMarker
        v-if="marker.data?.data.oneofKind === 'rectangle' && canMoveMarker(marker)"
        :lat-lng="[marker.data.data.rectangle.endY, marker.data.data.rectangle.endX]"
        :draggable="true"
        @dragstart="onMarkerDragStart"
        @dragend="onRectanglePointDragEnd($event, marker, false)"
    >
        <LIcon :icon-size="dragHandleSize" :icon-anchor="dragHandleAnchor">
            <div class="size-full rounded-sm border border-black bg-white/90"></div>
        </LIcon>
    </LMarker>

    <LPolygon
        v-if="marker.data?.data.oneofKind === 'polygon'"
        :name="marker.name"
        :lat-lngs="[
            [marker.y, marker.x],
            ...marker.data.data.polygon.points.map((point) => [point.y, point.x] satisfies LatLngTuple),
        ]"
        :color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-color="marker.color ?? livemap.markerMarkers.fallbackColor"
        :fill-opacity="(marker.data.data.polygon.opacity ?? 15) / 100"
        :options="{ markerMarker: marker }"
    >
        <MarkerMarkerPopup :marker="marker" />
    </LPolygon>
    <LMarker
        v-for="(point, idx) in canMoveMarker(marker) ? getPolygonDragPoints(marker) : []"
        :key="`poly-handle-${marker.id}-${idx}`"
        :lat-lng="[point.y, point.x]"
        :draggable="true"
        @dragstart="onMarkerDragStart"
        @dragend="onPolygonPointDragEnd($event, marker, idx)"
    >
        <LIcon :icon-size="dragHandleSize" :icon-anchor="dragHandleAnchor">
            <div class="size-full rounded-sm border border-black bg-white/90"></div>
        </LIcon>
    </LMarker>

    <LMarker
        v-if="marker.data?.data.oneofKind === 'icon'"
        :name="marker.name"
        :lat-lng="[marker.y, marker.x]"
        :draggable="canMoveMarker(marker)"
        :options="{ markerMarker: marker }"
        @click="$emit('selected')"
        @dragstart="onMarkerDragStart"
        @dragend="onMarkerDragEnd($event, marker)"
    >
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <component
                :is="resolveIconComponent(convertDynamicIconNameToComponent(marker.data?.data.icon.icon))"
                class="size-full"
                :style="{ color: marker.color ?? 'currentColor' }"
            />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>

    <LMarker
        v-if="
            marker.type !== MarkerType.DOT &&
            marker.data?.data.oneofKind !== 'circle' &&
            marker.data?.data.oneofKind !== 'rectangle' &&
            marker.data?.data.oneofKind !== 'polygon' &&
            marker.data?.data.oneofKind !== 'icon'
        "
        :name="marker.name"
        :lat-lng="[marker.y, marker.x]"
        :draggable="canMoveMarker(marker)"
        :options="{ markerMarker: marker }"
        @click="$emit('selected')"
        @dragstart="onMarkerDragStart"
        @dragend="onMarkerDragEnd($event, marker)"
    >
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <component :is="resolveIconComponent()" :fill="marker.color ?? 'currentColor'" />
        </LIcon>

        <MarkerMarkerPopup :marker="marker" />
    </LMarker>
</template>

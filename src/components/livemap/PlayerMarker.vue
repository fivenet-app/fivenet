<script lang="ts" setup>
import { LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import L from 'leaflet';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { User } from '~~/gen/ts/resources/users/users';

const props = withDefaults(
    defineProps<{
        marker: UserMarker;
        activeChar: null | User;
        size?: number;
    }>(),
    {
        size: 20,
    },
);

defineEmits<{
    (e: 'select'): void;
}>();

if (props.activeChar !== null && props.marker.user?.userId === props.activeChar.userId) {
    props.marker.marker!.color = 'FCAB10';
}

const iconAnchor: L.PointExpression = [props.size / 2, props.size];
const popupAnchor: L.PointExpression = [0, (props.size / 2) * -1];
const icon = new L.DivIcon({
    html: `<div>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -0.8 16 17.6" fill="${
            props.marker.marker?.color ? '#' + props.marker.marker?.color : 'currentColor'
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
    <LMarker
        ref="marker"
        :key="marker.marker!.id?.toString()"
        :latLng="[marker.marker!.y, marker.marker!.x]"
        :name="marker.marker!.name"
        :icon="icon"
        @click="$emit('select')"
        :z-index-offset="activeChar && marker.user?.identifier === activeChar.identifier ? 25 : 20"
    >
        <LPopup :options="{ closeButton: true }">
            <span class="font-semibold">{{ $t('common.employee', 2) }} {{ marker.user?.jobLabel }} </span>
            <br />
            <span class="italic">[{{ marker.user?.jobGrade }} {{ marker.user?.jobGradeLabel }}</span>
            <br />
            {{ marker.user?.firstname }} {{ marker.user?.lastname }}
        </LPopup>
    </LMarker>
</template>

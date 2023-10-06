<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import L from 'leaflet';
import { MapMarkerIcon } from 'mdi-vue3';
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
    (e: 'selected'): void;
}>();

if (props.activeChar !== null && props.marker.user?.userId === props.activeChar.userId) {
    props.marker.info!.color = 'FCAB10';
}

const inverseColor = computed(() => hexToRgb(props.marker.unit?.color ?? '#000000') ?? ({ r: 0, g: 0, b: 0 } as RGB));

const iconAnchor: L.PointExpression = [props.size / 2, props.size / 1];
const popupAnchor: L.PointExpression = [0.5, (props.size / 2) * -1];
</script>

<template>
    <LMarker
        :key="marker.info!.id?.toString()"
        :latLng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        @click="$emit('selected')"
        :z-index-offset="activeChar && marker.user?.identifier === activeChar.identifier ? 25 : 20"
    >
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="uppercase flex flex-col items-center dsp-status-error">
                <span
                    v-if="marker.unit"
                    class="rounded-md border-2 border-black/20 bg-clip-padding focus:outline-none inset-0 whitespace-nowrap"
                    :class="isColourBright(inverseColor) ? 'text-black' : 'text-white'"
                    :style="{ backgroundColor: '#' + props.marker.unit?.color ?? '000000' }"
                >
                    {{ marker.unit?.initials }}
                </span>
                <MapMarkerIcon class="w-full h-full" :style="{ color: '#' + props.marker.info?.color ?? '000000' }" />
            </div>
        </LIcon>
        <LPopup :options="{ closeButton: true }">
            <span class="font-semibold">{{ $t('common.employee', 2) }} {{ marker.user?.jobLabel }} </span>
            <ul role="list" class="flex flex-col">
                <li>
                    <span class="font-semibold"> {{ $t('common.name') }} </span>: {{ marker.user?.firstname }}
                    {{ marker.user?.lastname }}
                </li>
                <li>
                    <span class="font-semibold"> {{ $t('common.rank') }} </span>: {{ marker.user?.jobGradeLabel }} ({{
                        marker.user?.jobGrade
                    }})
                </li>
                <li v-if="marker.unit">
                    <span class="font-semibold">{{ $t('common.unit') }}</span
                    >: {{ marker.unit.name }} ({{ marker.unit.initials }})
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

<script lang="ts" setup>
import { LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import L from 'leaflet';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

const props = withDefaults(
    defineProps<{
        dispatch: Dispatch;
        size?: number;
    }>(),
    {
        size: 24,
    },
);

defineEmits<{
    (e: 'select'): void;
}>();

const iconClass = props.dispatch.status ? 'animate-dispatch' : '';
const iconAnchor: L.PointExpression | undefined = undefined;
const popupAnchor: L.PointExpression = [0, (props.size / 2) * -1];
const icon = new L.DivIcon({
    html: `<div class="${iconClass}">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#56db79" class="w-full h-full">
                <path d="M21,19V20H3V19L5,17V11C5,7.9 7.03,5.17 10,4.29C10,4.19 10,4.1 10,4A2,2 0 0,1 12,2A2,2 0 0,1 14,4C14,4.1 14,4.19 14,4.29C16.97,5.17 19,7.9 19,11V17L21,19M14,21A2,2 0 0,1 12,23A2,2 0 0,1 10,21"/>
            </svg>
        </div>`,
    iconSize: [props.size, props.size],
    iconAnchor,
    popupAnchor,
}) as L.Icon;
</script>

<template>
    <LMarker
        :key="dispatch.id?.toString()"
        :latLng="[dispatch.y, dispatch.x]"
        :name="dispatch.message"
        :z-index-offset="15"
        :icon="icon"
    >
        <LPopup :options="{ closeButton: true }">
            <IDCopyBadge class="mb-1" prefix="DSP" :id="dispatch.id" @action="$emit('select')" />
            <ul>
                <li>{{ $t('common.message') }}: {{ dispatch!.message }}</li>
                <li>{{ $t('common.description') }}: {{ dispatch!.description ?? 'N/A' }}</li>
                <li>{{ $t('common.sent_at') }}: {{ useLocaleTimeAgo(toDate(dispatch!.createdAt)!).value }}</li>
                <li class="italic">
                    {{ $t('common.sent_by') }}:
                    <span v-if="dispatch.anon">
                        {{ $t('common.anon') }}
                    </span>
                    <span v-else> {{ dispatch.user?.firstname }}, {{ dispatch.user?.lastname }} </span>
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

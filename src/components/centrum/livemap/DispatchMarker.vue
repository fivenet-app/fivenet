<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { BellIcon } from 'mdi-vue3';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { dispatchStatusToFillColor } from './helpers';

const props = withDefaults(
    defineProps<{
        dispatch: Dispatch;
        size?: number;
    }>(),
    {
        size: 24,
    },
);

const emits = defineEmits<{
    (e: 'select', dsp: Dispatch): void;
}>();

const animateStates = [
    DISPATCH_STATUS.NEW,
    DISPATCH_STATUS.UNIT_UNASSIGNED,
    DISPATCH_STATUS.UNASSIGNED,
    DISPATCH_STATUS.NEED_ASSISTANCE,
];
const status = props.dispatch.status?.status ?? 0;

const iconAnchor: L.PointExpression | undefined = undefined;
const popupAnchor: L.PointExpression = [0, (props.size / 2) * -1];

function selected(_: bigint | string) {
    emits('select', props.dispatch);
}
</script>

<template>
    <LMarker :key="dispatch.id?.toString()" :latLng="[dispatch.y, dispatch.x]" :name="dispatch.message" :z-index-offset="15">
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="uppercase flex flex-col items-center dsp-status-error">
                <span class="rounded-md bg-white border border-black">DSP-{{ props.dispatch.id }}</span>
                <BellIcon
                    class="w-full h-full"
                    :class="[animateStates.includes(status) ? 'animate-dispatch' : '', dispatchStatusToFillColor(status)]"
                />
                <span class="rounded-md bg-white border border-black">
                    {{ $t(`enums.centrum.DISPATCH_STATUS.${DISPATCH_STATUS[status]}`) }}</span
                >
            </div>
        </LIcon>
        <LPopup :options="{ closeButton: true }">
            <IDCopyBadge class="mb-1" prefix="DSP" :id="dispatch.id" :action="selected" />
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

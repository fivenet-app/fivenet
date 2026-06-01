<script lang="ts" setup>
import type { PointExpression } from 'leaflet';
import { defineAsyncComponent } from 'vue';
import { BellIcon } from 'mdi-vue3';
import { calculateDispatchZIndexOffset, dispatchStatusAnimate, dispatchStatusToFillColor } from '~/components/dispatch/helpers';
import LeafletLazyPopup from '~/components/livemap/LeafletLazyPopup.vue';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';

const props = withDefaults(
    defineProps<{
        dispatch: Dispatch;
        size?: number;
    }>(),
    {
        size: 22,
    },
);

const emit = defineEmits<{
    (e: 'selected', dsp: Dispatch): void;
}>();

const DispatchMarkerPopup = defineAsyncComponent(() => import('~/components/dispatch/livemap/DispatchMarkerPopup.vue'));

const iconAnchor: PointExpression = [props.size / 2, props.size * 1.65];
const popupAnchor: PointExpression = [0, -(props.size * 1.7)];

function selected() {
    emit('selected', props.dispatch);
}

const dispatchClasses = computed(() => [
    dispatchStatusToFillColor(props.dispatch.status?.status),
    dispatchStatusAnimate(props.dispatch.status?.status) ? 'animate-wiggle' : '',
]);

const zIndexOffset = computed(() => calculateDispatchZIndexOffset(props.dispatch.status?.status));
</script>

<template>
    <LMarker
        :key="dispatch.id"
        :lat-lng="[dispatch.y, dispatch.x]"
        :name="dispatch.id.toString()"
        :z-index-offset="zIndexOffset"
        :options="{ dispatchMarker: dispatch }"
    >
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="flex flex-col items-center uppercase">
                <span
                    class="inset-0 rounded-md border border-black/20 bg-neutral-50 bg-clip-padding whitespace-nowrap text-black"
                >
                    DSP-{{ dispatch.id }}
                </span>
                <BellIcon class="size-full" :class="dispatchClasses" />
            </div>
        </LIcon>

        <LeafletLazyPopup class="min-w-[175px] md:min-w-[305px]" :options="{ closeButton: false }">
            <template #default="{ popupOpen }">
                <DispatchMarkerPopup v-if="popupOpen" :dispatch="dispatch" @selected="selected" />
            </template>
        </LeafletLazyPopup>
    </LMarker>
</template>

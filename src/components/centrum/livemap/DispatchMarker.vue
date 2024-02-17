<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { type PointExpression } from 'leaflet';
import { BellIcon, CarEmergencyIcon } from 'mdi-vue3';
import { dispatchStatusAnimate, dispatchStatusToBGColor, dispatchStatusToFillColor } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';

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
    (e: 'goto', loc: Coordinate): void;
}>();

const iconAnchor: PointExpression = [props.size / 2, props.size * 1.65];
const popupAnchor: PointExpression = [0, -(props.size * 1.7)];

function selected(_: string | number | string) {
    emit('selected', props.dispatch);
}

const dispatchClasses = computed(() => [
    dispatchStatusToFillColor(props.dispatch.status?.status ?? 0),
    dispatchStatusAnimate(props.dispatch.status?.status ?? 0) ? 'animate-wiggle' : '',
]);
</script>

<template>
    <LMarker :key="dispatch.id" :lat-lng="[dispatch.y, dispatch.x]" :name="dispatch.message" :z-index-offset="10">
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="flex flex-col items-center uppercase">
                <span
                    class="inset-0 whitespace-nowrap rounded-md border-2 border-black/20 bg-neutral bg-clip-padding text-black hover:bg-[#f4f4f4] focus:outline-none"
                >
                    DSP-{{ dispatch.id }}
                </span>
                <BellIcon class="h-full w-full" :class="dispatchClasses" />
            </div>
        </LIcon>

        <LPopup :options="{ closeButton: true }">
            <div class="mb-1 flex items-center gap-2">
                <IDCopyBadge :id="dispatch.id" prefix="DSP" :action="selected" />
                <button
                    type="button"
                    :title="$t('common.delete')"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                    @click="selected(dispatch.id)"
                >
                    <CarEmergencyIcon class="h-5 w-5" />
                    <span class="ml-1">{{ $t('common.detail', 2) }}</span>
                </button>
            </div>

            <ul role="list" class="flex flex-col">
                <li>
                    <span class="font-semibold">{{ $t('common.sent_at') }}:</span> {{ $d(toDate(dispatch.createdAt), 'short') }}
                </li>
                <li class="inline-flex gap-1">
                    <span class="flex-initial">
                        <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                    </span>
                    <span class="flex-1">
                        <template v-if="dispatch.anon">
                            {{ $t('common.anon') }}
                        </template>
                        <CitizenInfoPopover v-else-if="dispatch.creator" :user="dispatch.creator" />
                        <template v-else>
                            {{ $t('common.unknown') }}
                        </template>
                    </span>
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.postal') }}:</span> {{ dispatch.postal ?? $t('common.na') }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.message') }}:</span> {{ dispatch.message }}
                </li>
                <li class="truncate">
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ dispatch.description ?? $t('common.na') }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.status') }}:</span>
                    <span class="ml-1" :class="dispatchStatusToBGColor(dispatch.status?.status)">
                        {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`) }}
                    </span>
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.attributes', 2) }}:</span>
                    <DispatchAttributes class="ml-1" :attributes="dispatch.attributes" />
                </li>
                <li class="inline-flex gap-1">
                    <span class="font-semibold">{{ $t('common.units') }}:</span>
                    <span v-if="dispatch.units.length === 0" class="italic">{{
                        $t('enums.centrum.StatusDispatch.UNASSIGNED')
                    }}</span>
                    <span v-else class="mr-1 grid grid-cols-2 gap-1">
                        <UnitInfoPopover
                            v-for="unit in dispatch.units"
                            :key="unit.unitId"
                            :unit="unit.unit"
                            :initials-only="true"
                            :badge="true"
                            :assignment="unit"
                        />
                    </span>
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

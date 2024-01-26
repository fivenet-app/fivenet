<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { type PointExpression } from 'leaflet';
import { BellIcon, CarEmergencyIcon } from 'mdi-vue3';
import { dispatchStatusAnimate, dispatchStatusToBGColor, dispatchStatusToFillColor } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';

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

const iconAnchor: PointExpression = [props.size / 2, props.size * 1.65];
const popupAnchor: PointExpression = [0, -(props.size * 1.7)];

function selected(_: string | number | string) {
    emit('selected', props.dispatch);
}

const dispatchClasses = computed(() => [
    dispatchStatusToFillColor(props.dispatch.status!.status),
    dispatchStatusAnimate(props.dispatch.status!.status) ? 'animate-wiggle' : '',
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
                    <span v-if="dispatch.units.length === 0" class="block">
                        {{ $t('common.unit', dispatch.units.length) }}
                    </span>
                    <div v-else class="rounded-md bg-base-800">
                        <ul role="list" class="divide-y divide-gray-200 text-sm font-medium">
                            <li
                                v-for="unit in dispatch.units"
                                :key="unit.unitId"
                                class="flex items-center justify-between py-3 pl-3 pr-4"
                            >
                                <div class="flex flex-1 items-center">
                                    <UnitInfoPopover
                                        :unit="unit.unit"
                                        :assignment="unit"
                                        class="flex items-center justify-center"
                                        text-class="text-gray-300"
                                    >
                                        <template #before>
                                            <AccountGroupIcon
                                                class="mr-1 h-5 w-5 flex-shrink-0 text-base-400"
                                                aria-hidden="true"
                                            />
                                        </template>
                                    </UnitInfoPopover>
                                    <span v-if="unit.expiresAt" class="ml-2 inline-flex flex-1 items-center truncate">
                                        -
                                        {{
                                            useLocaleTimeAgo(toDate(unit.expiresAt, timeCorrection), {
                                                showSecond: true,
                                                updateInterval: 1_000,
                                            }).value
                                        }}
                                    </span>
                                </div>
                            </li>
                        </ul>
                    </div>
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

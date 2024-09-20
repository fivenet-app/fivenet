<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import type { PointExpression } from 'leaflet';
import { MapMarkerIcon } from 'mdi-vue3';
import UnitDetailsSlideover from '~/components//centrum/units/UnitDetailsSlideover.vue';
import { unitStatusToBGColor } from '~/components/centrum/helpers';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import type { UserMarker } from '~~/gen/ts/resources/livemap/livemap';

const props = withDefaults(
    defineProps<{
        marker: UserMarker;
        size?: number;
        showUnitNames?: boolean;
        showUnitStatus?: boolean;
    }>(),
    {
        size: 20,
        showUnitNames: false,
        showUnitStatus: false,
    },
);

defineEmits<{
    (e: 'selected'): void;
}>();

const { livemap } = useAppConfig();

const slideover = useSlideover();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { goto } = useLivemapStore();

const centrumStore = useCentrumStore();
const { units } = storeToRefs(centrumStore);

const markerColor = computed(() => {
    if (activeChar.value !== null && props.marker.user?.userId === activeChar.value?.userId) {
        return livemap.userMarkers.activeCharColor;
    } else {
        return props.marker.info?.color ?? livemap.userMarkers.fallbackColor;
    }
});

const unit = computed(() => (props.marker.unitId !== undefined ? units.value.get(props.marker.unitId) : undefined));
const unitInverseColor = computed(() => {
    return hexToRgb(unit.value?.color ?? livemap.userMarkers.fallbackColor, RGBBlack)!;
});

const hasUnit = computed(() => props.showUnitNames && props.marker.unitId !== undefined);
const iconAnchor = computed<PointExpression | undefined>(() => [props.size / 2, props.size * (hasUnit.value ? 1.8 : 0.95)]);
const popupAnchor = computed<PointExpression>(() => (hasUnit.value ? [0, -(props.size * 1.7)] : [0, -(props.size * 0.8)]));

const unitStatusColor = computed(() => unitStatusToBGColor(unit.value?.status?.status ?? 0));
</script>

<template>
    <LMarker
        :key="`user_${marker.info!.id}`"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        :z-index-offset="activeChar === null || marker.user?.userId !== activeChar.userId ? 20 : 30"
        @click="$emit('selected')"
    >
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="flex flex-col items-center uppercase">
                <span
                    v-if="showUnitNames && unit"
                    class="inset-0 whitespace-nowrap rounded-md border border-black/20 bg-clip-padding"
                    :class="isColourBright(unitInverseColor) ? 'text-black' : 'text-neutral'"
                    :style="{ backgroundColor: unit?.color ?? livemap.userMarkers.fallbackColor }"
                >
                    {{ unit?.initials }}
                </span>
                <MapMarkerIcon class="size-full" :style="{ color: markerColor }" />
            </div>
            <div v-if="showUnitStatus && unit" class="pointer-events-none uppercase">
                <span class="absolute right-0 top-0 -mr-2 -mt-1.5 flex size-3">
                    <span
                        class="relative inset-0 inline-flex size-3 rounded-full border border-black/20"
                        :class="unitStatusColor"
                    ></span>
                </span>
            </div>
        </LIcon>

        <LPopup :options="{ closeButton: true }">
            <div class="flex flex-col gap-2">
                <div class="grid grid-cols-2 gap-2">
                    <UButton
                        v-if="marker.info?.x !== undefined && marker.info?.y !== undefined"
                        variant="link"
                        icon="i-mdi-map-marker"
                        :padded="false"
                        @click="goto({ x: marker.info?.x, y: marker.info?.y })"
                    >
                        <span class="truncate">
                            {{ $t('common.mark') }}
                        </span>
                    </UButton>

                    <UButton
                        v-if="can('CitizenStoreService.ListCitizens').value"
                        variant="link"
                        icon="i-mdi-account"
                        :padded="false"
                        :to="{ name: 'citizens-id', params: { id: marker.user?.userId ?? 0 } }"
                    >
                        <span class="truncate">
                            {{ $t('common.profile') }}
                        </span>
                    </UButton>

                    <UButton
                        v-if="can('JobsService.GetColleague').value && marker.user?.job === activeChar?.job"
                        variant="link"
                        icon="i-mdi-briefcase"
                        :padded="false"
                        :to="{ name: 'jobs-colleagues-id-info', params: { id: marker.user?.userId ?? 0 } }"
                    >
                        <span class="truncate">
                            {{ $t('common.colleague') }}
                        </span>
                    </UButton>

                    <PhoneNumberBlock
                        v-if="marker.user?.phoneNumber"
                        :number="marker.user?.phoneNumber"
                        :hide-number="true"
                        :show-label="true"
                        :padded="false"
                    />

                    <UButton
                        v-if="hasUnit && unit"
                        variant="link"
                        icon="i-mdi-group"
                        :padded="false"
                        @click="
                            slideover.open(UnitDetailsSlideover, {
                                unit: unit,
                            })
                        "
                    >
                        <span class="truncate">
                            {{ $t('common.unit') }}
                        </span>
                    </UButton>
                </div>

                <p class="inline-flex items-center gap-1">
                    <span class="font-semibold">{{ $t('common.employee', 2) }} {{ marker.user?.jobLabel }} </span>
                </p>

                <ul role="list">
                    <li>
                        <span class="font-semibold"> {{ $t('common.name') }} </span>: {{ marker.user?.firstname }}
                        {{ marker.user?.lastname }}
                    </li>
                    <li v-if="(marker.user?.jobGrade ?? 0) > 0 && marker.user?.jobGradeLabel">
                        <span class="font-semibold">{{ $t('common.rank') }}:</span> {{ marker.user?.jobGradeLabel }} ({{
                            marker.user?.jobGrade
                        }})
                    </li>
                    <li v-if="unit">
                        <span class="font-semibold">{{ $t('common.units') }}:</span> {{ unit.name }} ({{ unit.initials }})
                    </li>
                </ul>
            </div>
        </LPopup>
    </LMarker>
</template>

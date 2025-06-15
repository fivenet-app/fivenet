<script lang="ts" setup>
import type { LeafletMouseEvent, PointExpression } from 'leaflet';
import { MapMarkerIcon } from 'mdi-vue3';
import { unitStatusToBGColor } from '~/components/centrum/helpers';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import type { UserMarker } from '~~/gen/ts/resources/livemap/user_marker';
import UnitDetailsSlideover from '../centrum/units/UnitDetailsSlideover.vue';
import ColleagueName from '../jobs/colleagues/ColleagueName.vue';
import { checkIfCanAccessColleague } from '../jobs/colleagues/helpers';
import PhoneNumberBlock from '../partials/citizens/PhoneNumberBlock.vue';

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
    (e: 'selected', event: LeafletMouseEvent): void;
}>();

const { can, activeChar } = useAuth();

const { livemap } = useAppConfig();

const slideover = useSlideover();

const { goto } = useLivemapStore();

const centrumStore = useCentrumStore();
const { units } = storeToRefs(centrumStore);

const markerColor = computed(() => {
    if (activeChar.value !== null && props.marker.userId === activeChar.value?.userId) {
        return livemap.userMarkers.activeCharColor;
    } else {
        return props.marker.color ?? livemap.userMarkers.fallbackColor;
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

const markerRef = useTemplateRef('markerRef');
</script>

<template>
    <LMarker
        :key="`user_${marker.userId}`"
        ref="markerRef"
        :lat-lng="[marker.y, marker.x]"
        :z-index-offset="activeChar === null || marker.user?.userId !== activeChar.userId ? 20 : 30"
        @click="$emit('selected', $event)"
    >
        <LIcon
            :icon-anchor="iconAnchor"
            :popup-anchor="popupAnchor"
            :icon-size="[size, size]"
            :options="{
                pmIgnore: true,
            }"
        >
            <div class="flex flex-col items-center uppercase">
                <span
                    v-if="showUnitNames && unit"
                    class="inset-0 whitespace-nowrap rounded-md border border-black/20 bg-clip-padding"
                    :class="isColorBright(unitInverseColor) ? 'text-black' : 'text-neutral'"
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

        <LPopup class="min-w-[175px]" :options="{ closeButton: false }">
            <UCard
                class="-my-[13px] -ml-[20px] -mr-[24px] flex flex-col"
                :ui="{ body: { padding: 'px-2 py-2 sm:px-4 sm:p-2' } }"
            >
                <template #header>
                    <div class="grid grid-cols-2 gap-2">
                        <UButton
                            v-if="marker.x !== undefined && marker.y !== undefined"
                            variant="link"
                            icon="i-mdi-map-marker"
                            :padded="false"
                            block
                            @click="goto({ x: marker.x, y: marker.y })"
                        >
                            <span class="truncate">
                                {{ $t('common.mark') }}
                            </span>
                        </UButton>

                        <UButton
                            v-if="can('citizens.CitizensService.ListCitizens').value"
                            variant="link"
                            icon="i-mdi-account"
                            :padded="false"
                            block
                            :to="{ name: 'citizens-id', params: { id: marker.user?.userId ?? 0 } }"
                        >
                            <span class="truncate">
                                {{ $t('common.profile') }}
                            </span>
                        </UButton>

                        <UButton
                            v-if="
                                can('jobs.JobsService.GetColleague').value &&
                                marker.user &&
                                marker.user?.job === activeChar?.job &&
                                checkIfCanAccessColleague(marker.user, 'jobs.JobsService.GetColleague')
                            "
                            variant="link"
                            icon="i-mdi-briefcase"
                            :padded="false"
                            block
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
                            block
                        />

                        <UButton
                            v-if="hasUnit && unit"
                            variant="link"
                            icon="i-mdi-group"
                            :padded="false"
                            block
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
                </template>

                <p class="inline-flex items-center gap-1">
                    <span class="font-semibold">{{ $t('common.employee', 2) }} {{ marker.user?.jobLabel }} </span>
                </p>

                <ul role="list">
                    <li>
                        <span class="font-semibold"> {{ $t('common.name') }} </span>:
                        <ColleagueName v-if="marker.user" :colleague="marker.user" />
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
            </UCard>
        </LPopup>
    </LMarker>
</template>

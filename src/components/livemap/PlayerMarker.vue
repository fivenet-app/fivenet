<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { type PointExpression } from 'leaflet';
import { AccountIcon, GroupIcon, MapMarkerIcon } from 'mdi-vue3';
import UnitDetails from '~/components//centrum/units/UnitDetails.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { unitStatusToBGColor } from '~/components/centrum/helpers';
import type { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { useAuthStore } from '~/store/auth';
import { useCentrumStore } from '~/store/centrum';

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
    (e: 'goto', loc: Coordinate): void;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const centrumStore = useCentrumStore();
const { units } = storeToRefs(centrumStore);

function getMarkerColor(): string {
    if (activeChar !== null && props.marker.user?.userId === activeChar.value?.userId) {
        return '#fcab10';
    } else {
        return props.marker.info?.color ?? '#8d81f2';
    }
}

const unit = computed(() => (props.marker.unitId !== undefined ? units.value.get(props.marker.unitId) : undefined));
const unitInverseColor = computed(() => {
    return hexToRgb(unit.value?.color ?? '#8d81f2') ?? ({ r: 0, g: 0, b: 0 } as RGB);
});

const hasUnit = computed(() => props.showUnitNames && props.marker.unitId !== undefined);
const iconAnchor = computed<PointExpression | undefined>(() => [props.size / 2, props.size * (hasUnit.value ? 1.8 : 0.95)]);
const popupAnchor = computed<PointExpression>(() => (hasUnit.value ? [0, -(props.size * 1.7)] : [0, -(props.size * 0.8)]));

const unitStatusColor = computed(() => unitStatusToBGColor(unit.value?.status?.status ?? 0));

const openUnit = ref(false);
</script>

<template>
    <UnitDetails
        v-if="hasUnit && unit !== undefined"
        :unit="unit"
        :open="openUnit"
        @close="openUnit = false"
        @goto="$emit('goto', $event)"
    />

    <LMarker
        :key="`user_${marker.info!.id}`"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        :z-index-offset="activeChar === null || marker.user?.identifier !== activeChar.identifier ? 20 : 30"
        @click="$emit('selected')"
    >
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="flex flex-col items-center uppercase">
                <span
                    v-if="showUnitNames && unit"
                    class="inset-0 whitespace-nowrap rounded-md border-2 border-black/20 bg-clip-padding"
                    :class="isColourBright(unitInverseColor) ? 'text-black' : 'text-neutral'"
                    :style="{ backgroundColor: unit?.color ?? '#8d81f2' }"
                >
                    {{ unit?.initials }}
                </span>
                <MapMarkerIcon class="size-full" :style="{ color: getMarkerColor() }" aria-hidden="true" />
            </div>
            <div v-if="showUnitStatus && unit" class="pointer-events-none uppercase">
                <span class="absolute right-0 top-0 -mr-2 -mt-1.5 flex size-3">
                    <span
                        class="relative inset-0 inline-flex size-3 rounded-full border-2 border-black/20"
                        :class="unitStatusColor"
                    ></span>
                </span>
            </div>
        </LIcon>

        <LPopup :options="{ closeButton: true }">
            <div
                v-if="can('CitizenStoreService.ListCitizens') || marker.user?.phoneNumber || hasUnit"
                class="mb-1 flex items-center gap-2"
            >
                <UButton
                    v-if="marker.info?.x && marker.info?.y"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                    @click="$emit('goto', { x: marker.info?.x, y: marker.info?.y })"
                >
                    <MapMarkerIcon class="size-5" aria-hidden="true" />
                    <span class="ml-1">{{ $t('common.mark') }}</span>
                </UButton>
                <NuxtLink
                    v-if="can('CitizenStoreService.ListCitizens')"
                    :to="{ name: 'citizens-id', params: { id: marker.user?.userId ?? 0 } }"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                >
                    <AccountIcon class="size-5" aria-hidden="true" />
                    <span class="ml-1">{{ $t('common.profile') }}</span>
                </NuxtLink>
                <PhoneNumberBlock
                    v-if="marker.user?.phoneNumber"
                    :number="marker.user?.phoneNumber"
                    :hide-number="true"
                    :show-label="true"
                    width="w-4"
                />
                <UButton
                    v-if="hasUnit"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                    @click="openUnit = true"
                >
                    <GroupIcon class="size-4" aria-hidden="true" />
                    <span class="ml-1">
                        {{ $t('common.unit') }}
                    </span>
                </UButton>
            </div>
            <span class="font-semibold">{{ $t('common.employee', 2) }} {{ marker.user?.jobLabel }} </span>
            <ul role="list" class="flex flex-col">
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
                    <span class="font-semibold">{{ $t('common.unit') }}:</span> {{ unit.name }} ({{ unit.initials }})
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

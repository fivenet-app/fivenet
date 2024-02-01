<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { type PointExpression } from 'leaflet';
import { AccountIcon, GroupIcon, MapMarkerIcon } from 'mdi-vue3';
import UnitDetails from '~/components//centrum/units/UnitDetails.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { User } from '~~/gen/ts/resources/users/users';
import { unitStatusToBGColor } from '~/components/centrum/helpers';

const props = withDefaults(
    defineProps<{
        marker: UserMarker;
        activeChar: null | User;
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

function updateMarkerColor(): void {
    if (props.activeChar !== null && props.marker.user?.userId === props.activeChar.userId) {
        props.marker.info!.color = 'FCAB10';
    }
}

watch(props, () => updateMarkerColor());

updateMarkerColor();

const inverseColor = computed(() => hexToRgb(props.marker.unit?.color ?? '#8d81f2') ?? ({ r: 0, g: 0, b: 0 } as RGB));

const hasUnit = computed(() => props.showUnitNames && props.marker.unit !== undefined);
const iconAnchor = computed<PointExpression | undefined>(() => [props.size / 2, props.size * (hasUnit.value ? 1.8 : 0.95)]);
const popupAnchor = computed<PointExpression>(() => (hasUnit.value ? [0, -(props.size * 1.7)] : [0, -(props.size * 0.8)]));

const unitStatusColor = computed(() => unitStatusToBGColor(props.marker.unit?.status?.status ?? 0));

const openUnit = ref(false);
</script>

<template>
    <UnitDetails
        v-if="hasUnit && marker.unit !== undefined"
        :unit="marker.unit"
        :open="openUnit"
        @close="openUnit = false"
        @goto="$emit('goto', $event)"
    />

    <LMarker
        :key="marker.info!.id"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        :z-index-offset="activeChar && marker.user?.identifier === activeChar.identifier ? 30 : 20"
        @click="$emit('selected')"
    >
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="flex flex-col items-center uppercase">
                <span
                    v-if="showUnitNames && marker.unit"
                    class="inset-0 whitespace-nowrap rounded-md border-2 border-black/20 bg-clip-padding focus:outline-none"
                    :class="isColourBright(inverseColor) ? 'text-black' : 'text-neutral'"
                    :style="{ backgroundColor: '#' + (marker.unit?.color ?? '8d81f2') }"
                >
                    {{ marker.unit?.initials }}
                </span>
                <MapMarkerIcon class="h-full w-full" :style="{ color: '#' + (marker.info?.color ?? '8d81f2') }" />
            </div>
            <div v-if="showUnitStatus && marker.unit" class="pointer-events-none uppercase">
                <span class="absolute right-0 top-0 -mr-2 -mt-1.5 flex h-3 w-3">
                    <span
                        class="relative inset-0 inline-flex h-3 w-3 rounded-full border-2 border-black/20"
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
                <NuxtLink
                    v-if="can('CitizenStoreService.ListCitizens')"
                    :to="{ name: 'citizens-id', params: { id: marker.user?.userId ?? 0 } }"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                >
                    <AccountIcon class="h-5 w-5" />
                    <span class="ml-1">{{ $t('common.profile') }}</span>
                </NuxtLink>
                <PhoneNumberBlock
                    v-if="marker.user?.phoneNumber"
                    :number="marker.user?.phoneNumber"
                    :hide-number="true"
                    :show-label="true"
                    width="w-4"
                />
                <button
                    v-if="hasUnit"
                    type="button"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                    @click="openUnit = true"
                >
                    <GroupIcon class="h-4 w-4" />
                    <span class="ml-1">
                        {{ $t('common.unit') }}
                    </span>
                </button>
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
                <li v-if="marker.unit">
                    <span class="font-semibold">{{ $t('common.unit') }}:</span> {{ marker.unit.name }} ({{
                        marker.unit.initials
                    }})
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

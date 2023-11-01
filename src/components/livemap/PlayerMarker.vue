<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import L from 'leaflet';
import { AccountIcon, GroupIcon, MapMarkerIcon } from 'mdi-vue3';
import UnitDetails from '~/components//centrum/units/UnitDetails.vue';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';
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
}>();

function updateMarkerColor(): void {
    if (props.activeChar !== null && props.marker.user?.userId === props.activeChar.userId) {
        props.marker.info!.color = 'FCAB10';
    }
}

watch(props, () => updateMarkerColor());

updateMarkerColor();

const inverseColor = computed(() => hexToRgb(props.marker.unit?.color ?? '#000000') ?? ({ r: 0, g: 0, b: 0 } as RGB));

const hasUnit = computed(() => props.showUnitNames && props.marker.unit !== undefined);
const iconAnchor = computed<L.PointExpression | undefined>(() => [props.size / 2, props.size * (hasUnit.value ? 1.8 : 0.95)]);
const popupAnchor = computed<L.PointExpression>(() => (hasUnit.value ? [0, -(props.size * 1.7)] : [0, -(props.size * 0.8)]));

const unitStatusColor = computed(() => unitStatusToBGColor(props.marker.unit?.status?.status ?? 0));

const openUnit = ref(false);
</script>

<template>
    <UnitDetails
        v-if="hasUnit && props.marker.unit !== undefined"
        :unit="props.marker.unit"
        :open="openUnit"
        @close="openUnit = false"
    />

    <LMarker
        :key="marker.info!.id?.toString()"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        :z-index-offset="activeChar && marker.user?.identifier === activeChar.identifier ? 30 : 20"
        @click="$emit('selected')"
    >
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="uppercase flex flex-col items-center">
                <span
                    v-if="showUnitNames && marker.unit"
                    class="rounded-md border-2 border-black/20 bg-clip-padding focus:outline-none inset-0 whitespace-nowrap"
                    :class="isColourBright(inverseColor) ? 'text-black' : 'text-neutral'"
                    :style="{ backgroundColor: '#' + props.marker.unit?.color ?? '000000' }"
                >
                    {{ marker.unit?.initials }}
                </span>
                <MapMarkerIcon class="w-full h-full" :style="{ color: '#' + props.marker.info?.color ?? '000000' }" />
            </div>
            <div v-if="showUnitStatus && marker.unit" class="uppercase pointer-events-none">
                <span class="flex absolute h-3 w-3 top-0 right-0 -mt-1.5 -mr-2">
                    <span
                        class="relative inline-flex rounded-full h-3 w-3 border-2 border-black/20 inset-0"
                        :class="unitStatusColor"
                    ></span>
                </span>
            </div>
        </LIcon>
        <LPopup :options="{ closeButton: true }">
            <div class="flex items-center gap-2 mb-1">
                <NuxtLink
                    :to="{ name: 'citizens-id', params: { id: marker.user?.userId ?? 0 } }"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                >
                    <AccountIcon class="w-6 h-6" />
                    <span class="ml-1">{{ $t('common.profile') }}</span>
                </NuxtLink>
                <PhoneNumber
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
                    <GroupIcon class="w-4 h-4" />
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
                <li>
                    <span class="font-semibold">{{ $t('common.rank') }}</span
                    >: {{ marker.user?.jobGradeLabel }} ({{ marker.user?.jobGrade }})
                </li>
                <li v-if="marker.unit">
                    <span class="font-semibold">{{ $t('common.unit') }}</span
                    >: {{ marker.unit.name }} ({{ marker.unit.initials }})
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

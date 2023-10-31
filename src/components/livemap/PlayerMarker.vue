<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import L from 'leaflet';
import { AccountIcon, GroupIcon, MapMarkerIcon } from 'mdi-vue3';
import Details from '~/components//centrum/units/Details.vue';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { User } from '~~/gen/ts/resources/users/users';

const props = withDefaults(
    defineProps<{
        marker: UserMarker;
        activeChar: null | User;
        size?: number;
        showUnitNames?: boolean;
    }>(),
    {
        size: 20,
        showUnitNames: false,
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

const openUnit = ref(false);
</script>

<template>
    <Details
        v-if="hasUnit && props.marker.unit !== undefined"
        :unit="props.marker.unit"
        :open="openUnit"
        @close="openUnit = false"
    />

    <LMarker
        :key="marker.info!.id?.toString()"
        :latLng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        @click="$emit('selected')"
        :z-index-offset="activeChar && marker.user?.identifier === activeChar.identifier ? 30 : 20"
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
                    @click="openUnit = true"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
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

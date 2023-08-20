<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { LeafletMouseEvent } from 'leaflet';
import 'leaflet-contextmenu';
import 'leaflet-contextmenu/dist/leaflet.contextmenu.min.css';
import 'leaflet/dist/leaflet.css';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { GenericMarker } from '~~/gen/ts/resources/livemap/livemap';
import BaseMap from './BaseMap.vue';
import PlayerAndMarkersLayer from './PlayerAndMarkersLayer.vue';
import FollowSelectedMarker from './controls/FollowSelectedMarker.vue';
import MarkerResize from './controls/MarkerResize.vue';

withDefaults(
    defineProps<{
        centerSelectedMarker?: boolean;
        markerResize?: boolean;
        filterPlayers?: boolean;
    }>(),
    {
        centerSelectedMarker: true,
        markerResize: true,
        filterPlayers: true,
    },
);

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
const livemapStore = useLivemapStore();
const { error, abort, location } = storeToRefs(livemapStore);
const { startStream, stopStream } = livemapStore;

interface ContextmenuItem {
    text: string;
    callback: (e: LeafletMouseEvent) => void;
}

const mapOptions = {
    contextmenu: true,
    contextmenuWidth: 150,
    contextmenuItems: [] as ContextmenuItem[],
};

if (can('CentrumService.CreateDispatch')) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.create_dispatch.title'),
        callback: (e: LeafletMouseEvent) => {
            location.value = { x: e.latlng.lat, y: e.latlng.lng };
            createDispatchOpen.value = true;
        },
    });
}

onBeforeUnmount(() => {
    stopStream();
});

const createDispatchOpen = ref(false);

const selectedMarker = ref<GenericMarker | undefined>();

watch(selectedMarker, () => applySelectedMarkerCentering());

async function applySelectedMarkerCentering(): Promise<void> {
    if (selectedMarker.value === undefined) return;
    if (!livemap.value.centerSelectedMarker) return;

    location.value = { x: selectedMarker.value.x, y: selectedMarker.value.y };
}
</script>

<style>
.animate-dispatch {
    animation: wiggle 1s infinite;
}

@keyframes wiggle {
    0% {
        transform: rotate(0deg);
    }

    80% {
        transform: rotate(0deg);
    }

    85% {
        transform: rotate(5deg);
    }

    95% {
        transform: rotate(-5deg);
    }

    100% {
        transform: rotate(0deg);
    }
}
</style>

<template>
    <div class="relative w-full h-full z-0">
        <div
            v-if="error || (!error && abort === undefined)"
            class="absolute inset-0 flex justify-center items-center z-20 bg-gray-600/70"
        >
            <DataPendingBlock v-if="!error && abort === undefined" :message="$t('components.livemap.starting_datastream')" />
            <DataErrorBlock v-else-if="error" :title="$t('components.livemap.failed_datastream')" :retry="startStream" />
        </div>
        <BaseMap :map-options="mapOptions">
            <template v-slot:default>
                <PlayerAndMarkersLayer @marker-selected="selectedMarker = $event.marker" />

                <LControl position="bottomright" v-if="centerSelectedMarker">
                    <div class="form-control flex flex-col gap-2">
                        <MarkerResize />
                        <FollowSelectedMarker />
                    </div>
                </LControl>

                <slot />
            </template>
            <template v-slot:afterMap>
                <slot name="afterMap" />
            </template>
        </BaseMap>
    </div>
</template>

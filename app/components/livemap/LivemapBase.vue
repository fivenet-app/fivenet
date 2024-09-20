<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import type { ContextMenuItemClickEvent, type MapOptions } from 'leaflet';
import DispatchCreateOrUpdateSlideover from '~/components/centrum/dispatches/DispatchCreateOrUpdateSlideover.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import MapMarkersLayer from '~/components/livemap/MapMarkersLayer.vue';
import MapTempMarker from '~/components/livemap/MapTempMarker.vue';
import MapUsersLayer from '~/components/livemap/MapUsersLayer.vue';
import MarkerCreateOrUpdateSlideover from '~/components/livemap/MarkerCreateOrUpdateSlideover.vue';
import ReconnectingPopup from '~/components/livemap/ReconnectingPopup.vue';
import PostalSearch from '~/components/livemap/controls/PostalSearch.vue';
import SettingsButton from '~/components/livemap/controls/SettingsButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { isNUIAvailable, setWaypoint } from '~/composables/nui';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';

defineProps<{
    showUnitNames?: boolean;
    showUnitStatus?: boolean;
}>();

const { t } = useI18n();

const slideover = useSlideover();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const livemapStore = useLivemapStore();
const { startStream } = livemapStore;
const { error, abort, reconnecting, initiated, location, showLocationMarker, selectedMarker } = storeToRefs(livemapStore);

const centrumStore = useCentrumStore();
const { reconnecting: reconnectingCentrum } = storeToRefs(centrumStore);

const mapOptions = {
    zoomControl: false,
    contextmenu: true,
    contextmenuWidth: 150,
    contextmenuItems: [],
    scrollWheelZoom: 'center',
    markerZoomAnimation: true,
} as MapOptions;

if (can('CentrumService.CreateDispatch').value) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.create_dispatch.title'),
        callback: (e: ContextMenuItemClickEvent) => {
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            showLocationMarker.value = true;

            slideover.open(DispatchCreateOrUpdateSlideover, {
                location: { x: e.latlng.lng, y: e.latlng.lat },
                onClose: () => (showLocationMarker.value = false),
            });
        },
    });
}
if (can('LivemapperService.CreateOrUpdateMarker').value) {
    mapOptions.contextmenuItems.push({
        text: t('components.livemap.create_marker.title'),
        callback: (e: ContextMenuItemClickEvent) => {
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            showLocationMarker.value = true;

            slideover.open(MarkerCreateOrUpdateSlideover, {
                location: { x: e.latlng.lng, y: e.latlng.lat },
                onClose: () => (showLocationMarker.value = false),
            });
        },
    });
}
if (isNUIAvailable()) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.livemap.mark_on_gps'),
        callback: (e: ContextMenuItemClickEvent) => setWaypoint(e.latlng.lng, e.latlng.lat),
    });
}

function addActiveLayer(name: string): void {
    if (!livemap.value.activeLayers.includes(name)) {
        livemap.value.activeLayers.push(name);
    }
}

function removeActiveLayer(name: string): void {
    const idx = livemap.value.activeLayers.indexOf(name);
    if (idx > -1) {
        livemap.value.activeLayers.splice(idx, 1);
    }
}

const reconnectingDebounced = useDebounce(reconnecting, 500);
const reconnectionCentrumDebounced = useDebounce(reconnectingCentrum, 500);
</script>

<template>
    <div class="relative z-0 size-full">
        <div
            v-if="!!error || !initiated || (abort === undefined && !reconnecting)"
            class="absolute inset-0 z-20 flex items-center justify-center bg-gray-600/70"
        >
            <DataErrorBlock v-if="error" :title="$t('components.livemap.failed_datastream')" :retry="startStream" />
            <DataPendingBlock
                v-else-if="!initiated || (abort === undefined && !reconnectingDebounced)"
                :message="$t('components.livemap.starting_datastream')"
            />
        </div>

        <BaseMap
            :map-options="mapOptions"
            @overlayadd="addActiveLayer($event.name)"
            @overlayremove="removeActiveLayer($event.name)"
        >
            <template #default>
                <LControl position="bottomright">
                    <SettingsButton />
                </LControl>

                <template v-if="can('LivemapperService.Stream').value">
                    <MapUsersLayer
                        :show-unit-names="showUnitNames"
                        :show-unit-status="showUnitStatus"
                        @user-selected="selectedMarker = $event"
                    />
                    <MapMarkersLayer />
                </template>

                <MapTempMarker />

                <slot />

                <LControl position="bottomleft">
                    <PostalSearch />
                </LControl>
            </template>

            <template #afterMap>
                <ReconnectingPopup v-if="reconnectingDebounced || reconnectionCentrumDebounced" />

                <slot name="afterMap" />
            </template>
        </BaseMap>
    </div>
</template>

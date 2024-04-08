<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { type LeafletMouseEvent } from 'leaflet';
import DispatchCreateOrUpdateSlideover from '~/components/centrum/dispatches/DispatchCreateOrUpdateSlideover.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import CreateOrUpdateMarkerSlideover from '~/components/livemap/CreateOrUpdateMarkerSlideover.vue';
import MapTempMarker from '~/components/livemap/MapTempMarker.vue';
import MarkersLayer from '~/components/livemap/MarkersLayer.vue';
import PlayersLayer from '~/components/livemap/PlayersLayer.vue';
import ReconnectingPopup from '~/components/livemap/ReconnectingPopup.vue';
import PostalSearch from '~/components/livemap/controls/PostalSearch.vue';
import SettingsButton from '~/components/livemap/controls/SettingsButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { isNUIAvailable, setWaypoint } from '~/composables/nui';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { MarkerInfo } from '~~/gen/ts/resources/livemap/livemap';

defineProps<{
    showUnitNames?: boolean;
    showUnitStatus?: boolean;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { t } = useI18n();

const slideover = useSlideover();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const livemapStore = useLivemapStore();
const { error, abort, reconnecting, initiated, location, showLocationMarker } = storeToRefs(livemapStore);
const { startStream } = livemapStore;

const centrumStore = useCentrumStore();
const { reconnecting: reconnectingCentrum } = storeToRefs(centrumStore);

interface ContextmenuItem {
    text: string;
    callback: (e: LeafletMouseEvent) => void;
}

const mapOptions = {
    zoomControl: false,
    contextmenu: true,
    contextmenuWidth: 150,
    contextmenuItems: [] as ContextmenuItem[],
};

if (can('CentrumService.CreateDispatch')) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.create_dispatch.title'),
        callback: (e: LeafletMouseEvent) => {
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            showLocationMarker.value = true;

            slideover.open(DispatchCreateOrUpdateSlideover, {
                location: { x: e.latlng.lng, y: e.latlng.lat },
                onClose: () => (showLocationMarker.value = false),
            });
        },
    });
}
if (can('LivemapperService.CreateOrUpdateMarker')) {
    mapOptions.contextmenuItems.push({
        text: t('components.livemap.create_marker.title'),
        callback: (e: LeafletMouseEvent) => {
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            showLocationMarker.value = true;

            slideover.open(CreateOrUpdateMarkerSlideover, {
                location: { x: e.latlng.lng, y: e.latlng.lat },
                onClose: () => (showLocationMarker.value = false),
            });
        },
    });
}
if (isNUIAvailable()) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.livemap.mark_on_gps'),
        callback: (e: LeafletMouseEvent) => setWaypoint(e.latlng.lng, e.latlng.lat),
    });
}

const selectedUserMarker = ref<MarkerInfo | undefined>();

watch(selectedUserMarker, () => applySelectedMarkerCentering());

async function applySelectedMarkerCentering(): Promise<void> {
    if (selectedUserMarker.value === undefined || !livemap.value.centerSelectedMarker) {
        return;
    }

    location.value = { x: selectedUserMarker.value.x, y: selectedUserMarker.value.y };
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
            v-if="error !== undefined || !initiated || (abort === undefined && !reconnecting)"
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

                <template v-if="can('LivemapperService.Stream')">
                    <PlayersLayer
                        :show-unit-names="showUnitNames"
                        :show-unit-status="showUnitStatus"
                        @user-selected="selectedUserMarker = $event.info"
                        @goto="$emit('goto', $event)"
                    />
                    <MarkersLayer @goto="$emit('goto', $event)" />
                </template>

                <MapTempMarker />

                <slot />

                <LControl position="topleft">
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

<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { type LeafletMouseEvent } from 'leaflet';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { isNUIAvailable, setWaypoint } from '~/composables/nui';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { MarkerInfo } from '~~/gen/ts/resources/livemap/livemap';
import BaseMap from '~/components/livemap/BaseMap.vue';
import CreateOrUpdateMarkerModal from '~/components/livemap/CreateOrUpdateMarkerModal.vue';
import PlayersLayer from '~/components/livemap/PlayersLayer.vue';
import MarkersLayer from '~/components/livemap/MarkersLayer.vue';
import PostalSearch from '~/components/livemap/controls/PostalSearch.vue';
import SettingsButton from '~/components/livemap/controls/SettingsButton.vue';
import DispatchCreateOrUpdateSlideover from '~/components/centrum/dispatches/DispatchCreateOrUpdateSlideover.vue';
import MapTempMarker from '~/components/livemap/MapTempMarker.vue';
import ReconnectingPopup from '~/components/livemap/ReconnectingPopup.vue';
import { useCentrumStore } from '~/store/centrum';

defineProps<{
    showUnitNames?: boolean;
    showUnitStatus?: boolean;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { t } = useI18n();

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
            openCreateDispatch.value = true;
        },
    });
}
if (can('LivemapperService.CreateOrUpdateMarker')) {
    mapOptions.contextmenuItems.push({
        text: t('components.livemap.create_marker.title'),
        callback: (e: LeafletMouseEvent) => {
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            showLocationMarker.value = true;
            openCreateMarker.value = true;
        },
    });
}
if (isNUIAvailable()) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.livemap.mark_on_gps'),
        callback: (e: LeafletMouseEvent) => setWaypoint(e.latlng.lng, e.latlng.lat),
    });
}

const openCreateDispatch = ref(false);
const openCreateMarker = ref(false);

watch(openCreateDispatch, () => {
    if (openCreateDispatch.value) {
        showLocationMarker.value = false;
    }
});
watch(openCreateMarker, () => {
    if (openCreateMarker.value) {
        showLocationMarker.value = false;
    }
});

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
        <DispatchCreateOrUpdateSlideover
            v-if="can('CentrumService.CreateDispatch')"
            :open="openCreateDispatch"
            @close="openCreateDispatch = false"
        />
        <CreateOrUpdateMarkerModal
            v-if="can('LivemapperService.CreateOrUpdateMarker')"
            :open="openCreateMarker"
            @close="openCreateMarker = false"
        />

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

<script lang="ts" setup>
import type { ContextMenuItemClickEvent, MapOptions } from 'leaflet';
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
import { setWaypoint } from '~/composables/nui';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { useSettingsStore } from '~/stores/settings';

defineProps<{
    showUnitNames?: boolean;
    showUnitStatus?: boolean;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const livemapStore = useLivemapStore();
const { startStream } = livemapStore;
const { error, stopping: stoppingLivemap, initiated, location, showLocationMarker, selectedMarker } = storeToRefs(livemapStore);

const centrumStore = useCentrumStore();
const { stopping: stoppingCentrum } = storeToRefs(centrumStore);

const mapOptions = {
    zoomControl: false,
    contextmenu: true,
    contextmenuWidth: 150,
    contextmenuItems: [],
    scrollWheelZoom: 'center',
    markerZoomAnimation: true,
} as MapOptions;

if (can('centrum.CentrumService/CreateDispatch').value) {
    const dispatchCreateOrUpdateSlideover = overlay.create(DispatchCreateOrUpdateSlideover);

    mapOptions.contextmenuItems.push({
        text: t('components.centrum.create_dispatch.title'),
        callback: (e: ContextMenuItemClickEvent) => {
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            showLocationMarker.value = true;

            dispatchCreateOrUpdateSlideover.open({
                location: { x: e.latlng.lng, y: e.latlng.lat },
                onClose: () => (showLocationMarker.value = false),
            });
        },
    });
}
if (can('livemap.LivemapService/CreateOrUpdateMarker').value) {
    const markerCreateOrUpdateSlideover = overlay.create(MarkerCreateOrUpdateSlideover);

    mapOptions.contextmenuItems.push({
        text: t('components.livemap.create_marker.title'),
        callback: (e: ContextMenuItemClickEvent) => {
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            showLocationMarker.value = true;

            markerCreateOrUpdateSlideover.open({
                location: { x: e.latlng.lng, y: e.latlng.lat },
                onClose: () => (showLocationMarker.value = false),
            });
        },
    });
}
if (nuiEnabled.value) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.livemap.mark_on_gps'),
        callback: (e: ContextMenuItemClickEvent) => setWaypoint(e.latlng.lng, e.latlng.lat),
    });
}

const inititedDebounced = useDebounce(initiated, 750);
const stoppingLivemapDebounced = useDebounce(stoppingLivemap, 500);
const stoppingCentrumDebounced = useDebounce(stoppingCentrum, 500);
</script>

<template>
    <div class="relative size-full">
        <div v-if="error" class="absolute inset-0 z-20 flex items-center justify-center bg-neutral-600/70">
            <DataErrorBlock :title="$t('components.livemap.failed_datastream')" :error="error" :retry="startStream" />
        </div>

        <BaseMap :map-options="mapOptions">
            <template #default>
                <SettingsButton />

                <template v-if="can('livemap.LivemapService/Stream').value">
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
                <ReconnectingPopup
                    v-if="!inititedDebounced || stoppingLivemapDebounced || stoppingCentrumDebounced"
                    :label="
                        !inititedDebounced
                            ? $t('components.livemap.starting_datastream')
                            : $t('components.livemap.restarting_datastream')
                    "
                />

                <slot name="afterMap" />
            </template>
        </BaseMap>
    </div>
</template>

<script lang="ts" setup>
import type { MapOptions } from 'leaflet';
import DispatchCreateOrUpdateSlideover from '~/components/dispatch/dispatches/DispatchCreateOrUpdateSlideover.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import MarkerCreateOrUpdateSlideover from '~/components/livemap/markers/CreateOrUpdateSlideover.vue';
import MarkersLayer from '~/components/livemap/MarkersLayer.vue';
import ReconnectingPopup from '~/components/livemap/ReconnectingPopup.vue';
import TempMarker from '~/components/livemap/TempMarker.vue';
import UsersLayer from '~/components/livemap/UsersLayer.vue';
import PostalSearch from '~/components/livemap/controls/PostalSearch.vue';
import SettingsButton from '~/components/livemap/controls/SettingsButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { setWaypoint } from '~/composables/nui';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { useSettingsStore } from '~/stores/settings';
import type { LivemapContextMenuItem } from '~/types/livemap';
import type { Perms } from '~~/gen/ts/perms';

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
    scrollWheelZoom: 'center',
    markerZoomAnimation: true,
} as MapOptions;

const dispatchCreateOrUpdateSlideover = overlay.create(DispatchCreateOrUpdateSlideover);
const markerCreateOrUpdateSlideover = overlay.create(MarkerCreateOrUpdateSlideover);
const activeCreateOverlay = ref<'dispatch' | 'marker'>();
const pendingMarkerCreateLocation = ref<{ x: number; y: number }>();

const { start: startMarkerCreateOrUpdateOpenTimeout, stop: stopMarkerCreateOrUpdateOpenTimeout } = useTimeoutFn(
    () => {
        if (activeCreateOverlay.value !== 'marker' || !pendingMarkerCreateLocation.value) return;

        markerCreateOrUpdateSlideover
            .open({
                location: pendingMarkerCreateLocation.value,
                onClose: () => (showLocationMarker.value = false),
            })
            .finally(() => {
                pendingMarkerCreateLocation.value = undefined;

                if (activeCreateOverlay.value === 'marker') {
                    activeCreateOverlay.value = undefined;
                }
            });
    },
    0,
    { immediate: false },
);

function isCreateOverlayActive(): boolean {
    return (
        activeCreateOverlay.value !== undefined ||
        overlay.isOpen(dispatchCreateOrUpdateSlideover.id) ||
        overlay.isOpen(markerCreateOrUpdateSlideover.id)
    );
}

const contextMenuItems = computed<LivemapContextMenuItem[]>(() =>
    (
        [
            {
                label: t('components.dispatch.create_dispatch.title'),
                icon: 'i-mdi-car-emergency',
                permission: 'centrum.DispatchesService/CreateDispatch' as Perms,
                disabled: isCreateOverlayActive(),
                onSelect: (latlng) => {
                    if (isCreateOverlayActive()) return;

                    activeCreateOverlay.value = 'dispatch';
                    location.value = { x: latlng.lng, y: latlng.lat };
                    showLocationMarker.value = true;

                    dispatchCreateOrUpdateSlideover
                        .open({
                            location: { x: latlng.lng, y: latlng.lat },
                            onClose: () => (showLocationMarker.value = false),
                        })
                        .finally(() => {
                            if (activeCreateOverlay.value === 'dispatch') {
                                activeCreateOverlay.value = undefined;
                            }
                        });
                },
            },
            {
                label: t('components.livemap.create_marker.title'),
                icon: 'i-mdi-map-marker-outline',
                permission: 'livemap.LivemapService/CreateOrUpdateMarker' as Perms,
                disabled: isCreateOverlayActive(),
                onSelect: (latlng) => {
                    if (isCreateOverlayActive()) return;

                    const markerLocation = { x: latlng.lng, y: latlng.lat };
                    activeCreateOverlay.value = 'marker';
                    pendingMarkerCreateLocation.value = markerLocation;
                    location.value = markerLocation;
                    showLocationMarker.value = true;

                    startMarkerCreateOrUpdateOpenTimeout();
                },
            },
            ...((nuiEnabled.value
                ? [
                      {
                          type: 'separator' as const,
                      },
                      {
                          label: t('components.dispatch.livemap.mark_on_gps'),
                          icon: 'i-mdi-crosshairs-gps',
                          onSelect: (latlng) => setWaypoint(latlng.lng, latlng.lat),
                      },
                  ]
                : []) satisfies LivemapContextMenuItem[]),
        ] satisfies LivemapContextMenuItem[]
    ).filter((item) => item.permission === undefined || can(item.permission).value),
);

const inititedDebounced = useDebounce(initiated, 750);
const stoppingLivemapDebounced = useDebounce(stoppingLivemap, 500);
const stoppingCentrumDebounced = useDebounce(stoppingCentrum, 500);

onBeforeUnmount(() => stopMarkerCreateOrUpdateOpenTimeout());
</script>

<template>
    <div class="relative size-full">
        <div v-if="error" class="absolute inset-0 z-20 flex items-center justify-center bg-neutral-600/70">
            <DataErrorBlock :title="$t('components.livemap.failed_datastream')" :error="error" :retry="startStream" />
        </div>

        <BaseMap ref="baseMapRef" :map-options="mapOptions" :context-menu-items="contextMenuItems">
            <template #default>
                <SettingsButton />

                <template v-if="can('livemap.LivemapService/Stream').value">
                    <UsersLayer
                        :show-unit-names="showUnitNames"
                        :show-unit-status="showUnitStatus"
                        @user-selected="selectedMarker = $event"
                    />

                    <MarkersLayer />
                </template>

                <TempMarker />

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

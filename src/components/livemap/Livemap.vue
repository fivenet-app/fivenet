<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { LeafletMouseEvent } from 'leaflet';
import 'leaflet-contextmenu';
import 'leaflet-contextmenu/dist/leaflet.contextmenu.min.css';
import 'leaflet/dist/leaflet.css';
import CreateOrUpdateModal from '~/components/centrum/dispatches/CreateOrUpdateModal.vue';
import { checkForNUI, setWaypoint } from '~/components/centrum/nui';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { GenericMarker } from '~~/gen/ts/resources/livemap/livemap';
import BaseMap from './BaseMap.vue';
import PlayerAndMarkersLayer from './PlayerAndMarkersLayer.vue';
import PostalSearch from './controls/PostalSearch.vue';
import Settings from './controls/Settings.vue';

defineProps<{
    mapOptions?: Record<string, any>;
}>();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
const livemapStore = useLivemapStore();
const { error, abort, location } = storeToRefs(livemapStore);
const { startStream } = livemapStore;

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
            location.value = { x: e.latlng.lng, y: e.latlng.lat };
            openCreateDispatch.value = true;
        },
    });
}
if (checkForNUI()) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.livemap.mark_on_gps'),
        callback: (e: LeafletMouseEvent) => setWaypoint(e.latlng.lng, e.latlng.lat),
    });
}

const openCreateDispatch = ref(false);

const selectedMarker = ref<GenericMarker | undefined>();

watch(selectedMarker, () => applySelectedMarkerCentering());

async function applySelectedMarkerCentering(): Promise<void> {
    if (selectedMarker.value === undefined) return;
    if (!livemap.value.centerSelectedMarker) return;

    location.value = { x: selectedMarker.value.x, y: selectedMarker.value.y };
}
</script>

<template>
    <div class="relative w-full h-full z-0">
        <CreateOrUpdateModal :open="openCreateDispatch" @close="openCreateDispatch = false" />
        <div
            v-if="error || (!error && abort === undefined)"
            class="absolute inset-0 flex justify-center items-center z-20 bg-gray-600/70"
        >
            <DataErrorBlock v-if="error" :title="$t('components.livemap.failed_datastream')" :retry="startStream" />
            <DataPendingBlock v-else-if="abort === undefined" :message="$t('components.livemap.starting_datastream')" />
        </div>
        <BaseMap :map-options="mapOptions">
            <template v-slot:default>
                <LControl position="bottomright">
                    <div class="form-control flex flex-col gap-2">
                        <Settings />
                    </div>
                </LControl>

                <PlayerAndMarkersLayer @marker-selected="selectedMarker = $event.marker" />

                <slot />

                <LControl position="topleft">
                    <PostalSearch />
                </LControl>
            </template>

            <template v-slot:afterMap>
                <slot name="afterMap" />
            </template>
        </BaseMap>
    </div>
</template>

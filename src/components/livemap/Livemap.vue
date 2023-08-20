<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { LControl, LLayerGroup, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import L, { LeafletMouseEvent } from 'leaflet';
import 'leaflet-contextmenu';
import 'leaflet-contextmenu/dist/leaflet.contextmenu.min.css';
import 'leaflet/dist/leaflet.css';
import CreateOrUpdateModal from '~/components/centrum/dispatches/CreateOrUpdateModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useUserSettingsStore } from '~/store/usersettings';
import { Job } from '~~/gen/ts/resources/jobs/jobs';
import { DispatchMarker, UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';
import BaseMap from './BaseMap.vue';
import PlayerMarker from './PlayerMarker.vue';
import PostalSearch from './controls/PostalSearch.vue';

withDefaults(
    defineProps<{
        centerSelectedMarker?: boolean;
        markerResize?: boolean;
        filterEmployee?: boolean;
        filterDispatch?: boolean;
    }>(),
    {
        centerSelectedMarker: true,
        markerResize: true,
        filterEmployee: true,
        filterDispatch: true,
    },
);

const { $grpc, $loading } = useNuxtApp();
const { t } = useI18n();
const userSettingsStore = useUserSettingsStore();
const authStore = useAuthStore();

const { livemapCenterSelectedMarker, livemapMarkerSize } = storeToRefs(userSettingsStore);

$loading.start();

const { activeChar } = storeToRefs(authStore);

const abort = ref<AbortController | undefined>();
const error = ref<string | null>(null);

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
            goto({ x: e.latlng.lat, y: e.latlng.lng });
            createDispatchOpen.value = true;
        },
    });
}

const markerDispatches = ref<Job[]>([]);
const markerPlayers = ref<Job[]>([]);
const selectedMarker = ref<bigint>();

watch(livemapCenterSelectedMarker, () => {
    applySelectedMarkerCentering();
});

const playerQuery = ref<string>('');
let playerMarkers: UserMarker[] = [];
const playerMarkersFiltered = computed(() =>
    playerMarkers.filter((m) => (m.user?.firstname + ' ' + m.user?.lastname).includes(playerQuery.value)),
);

const dispatchQuery = ref<string>('');
let dispatchMarkers: DispatchMarker[] = [];
const dispatchMarkersFiltered = computed(() =>
    dispatchMarkers.filter(
        (m) => m.marker?.popup?.includes(dispatchQuery.value) || m.marker?.name.includes(dispatchQuery.value),
    ),
);

let map: L.Map | undefined = undefined;

const location = ref<{ x: number; y: number }>({ x: 0, y: 0 });
defineExpose({
    location,
});

watch(location, () => goto({ x: location.value?.x, y: location.value?.y }));

async function startStream(): Promise<void> {
    if (abort.value !== undefined) return;

    console.debug('Livemap: Starting Data Stream');
    try {
        abort.value = new AbortController();

        const call = new LivemapperServiceClient($grpc.getTransport()).stream(
            {},
            {
                abort: abort.value.signal,
            },
        );

        for await (let resp of call.responses) {
            error.value = null;

            if (resp === undefined || !resp.jobsDispatches || !resp.jobsDispatches) {
                continue;
            }

            markerDispatches.value = resp.jobsDispatches;
            markerPlayers.value = resp.jobsUsers;

            playerMarkers = resp.users;
            dispatchMarkers = resp.dispatches;

            applySelectedMarkerCentering();
        }
    } catch (e) {
        const err = e as RpcError;
        error.value = err.message;
        // TODO Restart stream automatically if timeout occurs
        $loading.errored();
        stopStream();
    }

    console.debug('Livemap: Data Stream Ended');
}

async function stopStream(): Promise<void> {
    console.debug('Livemap: Stopping Data Stream');
    abort.value?.abort();
    abort.value = undefined;
}

async function applySelectedMarkerCentering(): Promise<void> {
    if (!livemapCenterSelectedMarker.value) return;
    if (selectedMarker.value === undefined) return;

    const marker =
        playerMarkers.find((m) => m.marker?.id === selectedMarker.value) ||
        playerMarkers.find((m) => m.marker?.id === selectedMarker.value);
    if (!marker) {
        selectedMarker.value = undefined;
        return;
    }
    if (!marker.marker) {
        selectedMarker.value = undefined;
        return;
    }

    goto({ x: marker.marker.x, y: marker.marker.y });
}

onBeforeUnmount(() => {
    stopStream();
    map = undefined;
});

async function setSelectedMarker(id: bigint): Promise<void> {
    setTimeout(() => {
        selectedMarker.value = id;
        applySelectedMarkerCentering();
    }, 100);
}

const createOrUpdateModal = ref<InstanceType<typeof CreateOrUpdateModal>>();
const createDispatchOpen = ref(false);

function goto(e: { x: number; y: number }) {
    if (createOrUpdateModal.value) {
        createOrUpdateModal.value.location = { x: e.x, y: e.y };
    }

    if (map) {
        location.value = { x: e.x, y: e.y };
    }
}

function onMapReady(map: L.Map): void {
    startStream();
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
            v-if="error || abort === undefined"
            class="absolute inset-0 flex justify-center items-center z-20"
            style="background-color: rgba(62, 60, 62, 0.5)"
        >
            <DataPendingBlock v-if="!error" :message="$t('components.livemap.starting_datastream')" />
            <DataErrorBlock v-else-if="error" :title="$t('components.livemap.failed_datastream')" :retry="startStream" />
            <DataPendingBlock v-else-if="!error" :message="$t('components.livemap.paused_datastream')" :paused="true" />
        </div>
        <BaseMap :map-options="mapOptions" @map-ready="onMapReady">
            <template v-slot:default>
                <LLayerGroup
                    v-for="job in markerPlayers"
                    :key="job.name"
                    :name="`${$t('common.employee', 2)} ${job.label}`"
                    layer-type="overlay"
                    :visible="true"
                >
                    <PlayerMarker
                        v-for="marker in playerMarkersFiltered.filter((p) => p.user?.job === job.name)"
                        :key="marker.marker!.id?.toString()"
                        :marker="marker"
                        :active-char="activeChar"
                        @selected="setSelectedMarker(marker.marker!.id)"
                    />
                </LLayerGroup>

                <LLayerGroup
                    v-for="job in markerDispatches"
                    :key="job.name"
                    :name="`${$t('common.dispatch', 2)} ${job.label}`"
                    layer-type="overlay"
                    :visible="true"
                >
                    <LMarker
                        v-for="marker in dispatchMarkersFiltered.filter((m) => m.job === job.name)"
                        :key="marker.marker!.id?.toString()"
                        :latLng="[marker.marker!.y, marker.marker!.x]"
                        :name="marker.marker!.name"
                        @click="setSelectedMarker(marker.marker!.id)"
                        :z-index-offset="marker.active ? 15 : 10"
                    >
                        <LPopup
                            :options="{ closeButton: false }"
                            :content="`<span class='font-semibold'>${$t('common.dispatch', 2)} ${marker.jobLabel}</span><br>${
                                marker.marker!.popup
                            }<br><span>${
                                useLocaleTimeAgo(toDate(marker.marker!.updatedAt)!).value
                            }</span><br><span class='italic'>${$t('common.sent_by')} ${marker.marker!.name}</span>`"
                        >
                        </LPopup>
                    </LMarker>
                </LLayerGroup>

                <LControl position="topleft" v-if="filterDispatch || filterEmployee">
                    <div class="form-control flex flex-col gap-2">
                        <div v-if="filterEmployee">
                            <input
                                v-model="playerQuery"
                                class="w-full"
                                type="text"
                                name="searchPlayer"
                                :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
                            />
                        </div>
                        <div v-if="filterDispatch">
                            <input
                                v-model="dispatchQuery"
                                class="w-full"
                                type="text"
                                name="searchDispatch"
                                :placeholder="`${$t('common.dispatch', 1)} ${$t('common.filter')}`"
                            />
                        </div>
                        <PostalSearch @goto="goto($event)" />
                    </div>
                </LControl>
                <LControl position="bottomright" v-if="centerSelectedMarker || markerResize">
                    <div class="form-control flex flex-col gap-2">
                        <div
                            v-if="centerSelectedMarker"
                            class="p-2 bg-neutral border border-[#6b7280] flex flex-row justify-center"
                        >
                            <span class="text-lg mr-2 text-[#6f7683]">{{
                                $t('components.livemap.center_selected_marker')
                            }}</span>
                            <input
                                v-model="livemapCenterSelectedMarker"
                                class="my-auto"
                                name="livemapCenterSelectedMarker"
                                type="checkbox"
                            />
                        </div>
                        <div v-if="markerResize" class="p-2 bg-neutral border border-[#6b7280] flex flex-row justify-center">
                            <span class="text-lg mr-2 text-[#6f7683]">{{ livemapMarkerSize }}</span>
                            <input
                                name="livemapMarkerSize"
                                type="range"
                                class="h-1.5 w-full cursor-grab rounded-full my-auto"
                                min="14"
                                max="34"
                                step="2"
                                :value="livemapMarkerSize"
                                @change="livemapMarkerSize = ($event.target as any).value"
                            />
                        </div>
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

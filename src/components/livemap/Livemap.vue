<script lang="ts" setup>
import { Combobox, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { LControl, LControlLayers, LLayerGroup, LMap, LMarker, LPopup, LTileLayer } from '@vue-leaflet/vue-leaflet';
import { watchDebounced } from '@vueuse/core';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import LoadingBar from '~/components/partials/LoadingBar.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import { useUserSettingsStore } from '~/store/usersettings';
import { ValueOf } from '~/utils/types';
import { Job } from '~~/gen/ts/resources/jobs/jobs';
import { DispatchMarker, UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';

const { $grpc, $loading } = useNuxtApp();
const userSettingsStore = useUserSettingsStore();
const authStore = useAuthStore();
const notifications = useNotificationsStore();
const route = useRoute();

const { livemapCenterSelectedMarker, livemapMarkerSize } = storeToRefs(userSettingsStore);

const { t } = useI18n();

$loading.start();

const { activeChar } = storeToRefs(authStore);

const abort = ref<AbortController | undefined>();
const error = ref<string | null>(null);

const centerX = 117.3;
const centerY = 172.8;
const scaleX = 0.02072;
const scaleY = 0.0205;

const customCRS = L.extend({}, L.CRS.Simple, {
    projection: L.Projection.LonLat,
    scale: function (zoom: number): number {
        return Math.pow(2, zoom);
    },
    zoom: function (sc: number): number {
        return Math.log(sc) / 0.6931471805599453;
    },
    distance: function (pos1: L.LatLng, pos2: L.LatLng): number {
        var x_difference = pos2.lng - pos1.lng;
        var y_difference = pos2.lat - pos1.lat;
        return Math.sqrt(x_difference * x_difference + y_difference * y_difference);
    },
    transformation: new L.Transformation(scaleX, centerX, -scaleY, centerY),
    infinite: true,
});

const backgroundColorList = {
    Atlas: '#0fa8d2',
    Satelite: '#143d6b',
    Road: '#1862ad',
    Postal: '#74aace',
} as const;
const backgroundColor = ref<ValueOf<typeof backgroundColorList>>(backgroundColorList.Postal);

const zoom = ref(2);
let center: L.PointExpression = [0, 0];
const attribution = '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>';

const markerDispatches = ref<Job[]>([]);
const markerPlayers = ref<Job[]>([]);
const selectedMarker = ref<number>();

watch(livemapCenterSelectedMarker, () => {
    applySelectedMarkerCentering();
});

const playerQuery = ref<string>('');
let playerMarkers: UserMarker[] = [];
const playerMarkersFiltered = ref<UserMarker[]>([]);

const dispatchQuery = ref<string>('');
let dispatchMarkers: DispatchMarker[] = [];
const dispatchMarkersFiltered = ref<DispatchMarker[]>([]);

async function applyPlayerQuery(): Promise<void> {
    playerMarkersFiltered.value = playerMarkers.filter((m) =>
        (m.user?.firstname + ' ' + m.user?.lastname).includes(playerQuery.value)
    );
}
async function applyDispatchQuery(): Promise<void> {
    dispatchMarkersFiltered.value = dispatchMarkers.filter(
        (m) => m.popup.includes(dispatchQuery.value) || m.name.includes(dispatchQuery.value)
    );
}

watchDebounced(
    playerQuery,
    async () => {
        applyPlayerQuery();
    },
    { debounce: 600, maxWait: 1750 }
);
watchDebounced(
    dispatchQuery,
    async () => {
        applyDispatchQuery();
    },
    { debounce: 600, maxWait: 1750 }
);

const mouseLat = ref<string>((0).toFixed(3));
const mouseLong = ref<string>((0).toFixed(3));

const currentHash = ref<string>('');
const isMoving = ref<boolean>(false);

let map: L.Map | undefined = undefined;

watch(currentHash, () => {
    window.location.replace(currentHash.value);
});

watchDebounced(
    isMoving,
    () => {
        if (!map || isMoving.value) return;

        const newHash = stringifyHash(map.getZoom(), map.getCenter().lat, map.getCenter().lng);
        if (currentHash.value !== newHash) currentHash.value = newHash;
    },
    { debounce: 1000, maxWait: 3000 }
);

async function updateBackground(layer: string): Promise<void> {
    switch (layer) {
        case 'Atlas':
            backgroundColor.value = backgroundColorList.Atlas;
            return;
        case 'Satelite':
            backgroundColor.value = backgroundColorList.Satelite;
            return;
        case 'Road':
            backgroundColor.value = backgroundColorList.Road;
            return;
        case 'Postal':
            backgroundColor.value = backgroundColorList.Postal;
            return;
    }
}

function stringifyHash(currZoom: number, centerLat: number, centerLong: number): string {
    const precision = Math.max(0, Math.ceil(Math.log(zoom.value) / Math.LN2));

    const hash = '#' + [currZoom, centerLat.toFixed(precision), centerLong.toFixed(precision)].join('/');
    return hash;
}

function parseHash(hash: string): { latlng: L.LatLng; zoom: number } | undefined {
    if (hash.indexOf('#') === 0) hash = hash.substring(1);

    const args = hash.split('/');
    if (args.length !== 3) return;

    const zoom = parseInt(args[0], 10);
    const lat = parseFloat(args[1]);
    const lng = parseFloat(args[2]);

    if (isNaN(zoom) || isNaN(lat) || isNaN(lng)) return;

    return {
        latlng: new L.LatLng(lat, lng),
        zoom,
    };
}

async function onMapReady($event: any): Promise<void> {
    map = $event as L.Map;

    const startingHash = route.hash;
    const startPos = parseHash(startingHash);
    if (startPos) $event.setView(startPos.latlng, startPos.zoom);

    map.on('baselayerchange', async (event: L.LayersControlEvent) => {
        updateBackground(event.name);
    });

    map.addEventListener('mousemove', async (event: L.LeafletMouseEvent) => {
        mouseLat.value = (Math.round(event.latlng.lat * 100000) / 100000).toFixed(3);
        mouseLong.value = (Math.round(event.latlng.lng * 100000) / 100000).toFixed(3);
    });

    map.on('movestart', async () => {
        isMoving.value = true;
    });

    map.on('moveend', async () => {
        isMoving.value = false;
    });

    setTimeout(() => {
        $loading.finish();
    }, 500);

    startDataStream();
}

async function startDataStream(): Promise<void> {
    if (abort.value !== undefined) return;

    console.debug('Livemap: Starting Data Stream');
    try {
        abort.value = new AbortController();

        const call = new LivemapperServiceClient($grpc.getTransport()).stream(
            {},
            {
                abort: abort.value.signal,
            }
        );

        for await (let resp of call.responses) {
            error.value = null;

            markerDispatches.value = resp.jobsDispatches;
            markerPlayers.value = resp.jobsUsers;

            playerMarkers = resp.users;
            dispatchMarkers = resp.dispatches;

            await applyPlayerQuery();
            await applyDispatchQuery();
            applySelectedMarkerCentering();
        }
    } catch (e) {
        const err = e as RpcError;
        error.value = err.message;
        $loading.errored();
        stopDataStream();
    }

    console.debug('Livemap: Data Stream Ended');
}

async function stopDataStream(): Promise<void> {
    console.debug('Livemap: Stopping Data Stream');
    abort.value?.abort();
    abort.value = undefined;
}

async function applySelectedMarkerCentering(): Promise<void> {
    if (!livemapCenterSelectedMarker.value) return;
    if (selectedMarker.value === undefined) return;

    const marker =
        playerMarkers.find((m) => m.id === selectedMarker.value) || playerMarkers.find((m) => m.id === selectedMarker.value);
    if (!marker) {
        selectedMarker.value = undefined;
        return;
    }

    map?.panTo([marker.y, marker.x], {
        animate: true,
        duration: 0.85,
    });
}

type TMarker<TType> = TType extends 'player' ? UserMarker : TType extends 'dispatch' ? DispatchMarker : never;

function getIcon<TType extends 'player' | 'dispatch'>(type: TType, marker: TMarker<TType>): L.DivIcon {
    let html = '';
    let color = marker.iconColor;
    let iconClass = '';
    let iconAnchor: L.PointExpression | undefined = undefined;
    let popupAnchor: L.PointExpression = [0, (livemapMarkerSize.value / 2) * -1];

    switch (type) {
        case 'player':
            {
                if (activeChar.value && (marker as UserMarker).user?.identifier === activeChar.value.identifier)
                    color = 'FCAB10';
                iconAnchor = [livemapMarkerSize.value / 2, livemapMarkerSize.value];
                popupAnchor = [0, livemapMarkerSize.value * -1];

                html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -0.8 16 17.6" fill="${
                    color ? '#' + color : 'currentColor'
                }" class="w-full h-full">
                    <path d="M8 16s6-5.686 6-10A6 6 0 0 0 2 6c0 4.314 6 10 6 10zm0-7a3 3 0 1 1 0-6 3 3 0 0 1 0 6z"/>
                </svg>`;
            }
            break;

        case 'dispatch':
            {
                if ((marker as DispatchMarker).active) iconClass = 'animate-dispatch';

                html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -0.8 16 17.6" fill="${
                    color ? '#' + color : 'currentColor'
                }" class="w-full h-full">
                    <path d="M8 16a2 2 0 0 0 2-2H6a2 2 0 0 0 2 2zm.995-14.901a1 1 0 1 0-1.99 0A5.002 5.002 0 0 0 3 6c0 1.098-.5 6-2 7h14c-1.5-1-2-5.902-2-7 0-2.42-1.72-4.44-4.005-4.901z"/>
                </svg>`;
            }
            break;
    }

    return new L.DivIcon({
        html: `<div class="${iconClass}">` + html + '</div>',
        iconSize: [livemapMarkerSize.value, livemapMarkerSize.value],
        iconAnchor,
        popupAnchor,
    });
}

onBeforeUnmount(() => {
    stopDataStream();
    map = undefined;
});

type Postal = {
    x: number;
    y: number;
    code: string;
};

const selectedPostal = ref<Postal | undefined>();
const postalQuery = ref('');
let postalsLoaded = false;
const postals = ref<Array<Postal>>([]);
const filteredPostals = ref<Array<Postal>>([]);

async function loadPostals(): Promise<void> {
    if (postalsLoaded) {
        return;
    }
    postalsLoaded = true;

    try {
        const response = await fetch('/data/postals.json');
        postals.value.push(...((await response.json()) as Postal[]));
    } catch (_) {
        notifications.dispatchNotification({
            title: t('notifications.failed_loading_postals.title'),
            content: t('notifications.failed_loading_postals.content'),
            type: 'error',
        });
        postalsLoaded = false;
    }
}

async function findPostal(): Promise<void> {
    if (postalQuery.value === '') {
        return;
    }

    let results = 0;
    filteredPostals.value.length = 0;
    filteredPostals.value = postals.value.filter((p) => {
        if (results >= 10) {
            return false;
        }
        const result = p.code.startsWith(postalQuery.value!);
        if (result) results++;
        return result;
    });
    if (filteredPostals.value.length === 0) {
        return;
    }
}

async function setSelectedMarker(id: number): Promise<void> {
    setTimeout(() => {
        selectedMarker.value = id;
        applySelectedMarkerCentering();
    }, 100);
}

watch(selectedPostal, () => {
    if (!selectedPostal.value) {
        return;
    }

    map?.flyTo([selectedPostal.value.y, selectedPostal.value.x], 5, {
        animate: true,
        duration: 0.85,
    });
});

watchDebounced(postalQuery, () => findPostal(), {
    debounce: 250,
    maxWait: 850,
});
</script>

<style>
.leaflet-div-icon {
    background: none;
    border: none;
}

.leaflet-div-icon svg path {
    stroke: #000;
    stroke-width: 0.75px;
    stroke-linejoin: round;
}

.leaflet-marker-icon {
    transition: transform 1s ease;
}

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
    <LoadingBar />
    <div class="relative w-full h-full z-0">
        <div
            v-if="error || abort === undefined"
            class="absolute inset-0 flex justify-center items-center z-20"
            style="background-color: rgba(62, 60, 62, 0.5)"
        >
            <DataPendingBlock v-if="!error" :message="$t('components.livemap.starting_datastream')" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('components.livemap.failed_datastream')"
                :retry="
                    () => {
                        startDataStream();
                    }
                "
            />
            <DataPendingBlock v-else-if="!error" :message="$t('components.livemap.paused_datastream')" :paused="true" />
        </div>

        <LMap
            class="z-0"
            v-model:zoom="zoom"
            v-model:center="center"
            :crs="customCRS"
            :min-zoom="1"
            :max-zoom="6"
            @click="selectedMarker = undefined"
            :inertia="false"
            :style="{ backgroundColor }"
            @ready="onMapReady($event)"
            :use-global-leaflet="false"
        >
            <LTileLayer
                url="/images/livemap/tiles/postal/{z}/{x}/{y}.png"
                layer-type="base"
                name="Postal"
                :no-wrap="true"
                :tms="true"
                :visible="true"
                :attribution="attribution"
            />
            <LTileLayer
                url="/images/livemap/tiles/atlas/{z}/{x}/{y}.png"
                layer-type="base"
                name="Atlas"
                :no-wrap="true"
                :tms="true"
                :visible="false"
                :attribution="attribution"
            />
            <LTileLayer
                url="/images/livemap/tiles/road/{z}/{x}/{y}.png"
                layer-type="base"
                name="Road"
                :no-wrap="true"
                :tms="true"
                :visible="false"
                :attribution="attribution"
            />
            <LTileLayer
                url="/images/livemap/tiles/satelite/{z}/{x}/{y}.png"
                layer-type="base"
                name="Satelite"
                :no-wrap="true"
                :tms="true"
                :visible="false"
                :attribution="attribution"
            />

            <LControlLayers />

            <LLayerGroup
                v-for="job in markerPlayers"
                :key="job.name"
                :name="`${$t('common.employee', 2)} ${job.label}`"
                layer-type="overlay"
                :visible="true"
            >
                <LMarker
                    v-for="marker in playerMarkersFiltered.filter((p) => p.user?.job === job.name)"
                    :key="marker.id"
                    :latLng="[marker.y, marker.x]"
                    :name="marker.name"
                    :icon="getIcon('player', marker) as L.Icon"
                    @click="setSelectedMarker(marker.id)"
                    :z-index-offset="activeChar && marker.user?.identifier === activeChar.identifier ? 25 : 20"
                >
                    <LPopup
                        :options="{ closeButton: false }"
                        :content="`<span class='font-semibold'>${$t('common.employee', 2)} ${
                            marker.user?.jobLabel
                        }</span><br><span class='italic'>[${marker.user?.jobGrade}] ${marker.user?.jobGradeLabel}</span><br>${
                            marker.user?.firstname
                        } ${marker.user?.lastname}`"
                    >
                    </LPopup>
                </LMarker>
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
                    :key="marker.id"
                    :latLng="[marker.y, marker.x]"
                    :name="marker.name"
                    :icon="getIcon('dispatch', marker) as L.Icon"
                    @click="setSelectedMarker(marker.id)"
                    :z-index-offset="marker.active ? 15 : 10"
                >
                    <LPopup
                        :options="{ closeButton: false }"
                        :content="`<span class='font-semibold'>${$t('common.dispatch', 2)} ${marker.jobLabel}</span><br>${marker.popup}<br><span>${useLocaleTimeAgo(toDate(marker.updatedAt)!).value}</span><br><span class='italic'>${$t('components.livemap.sent_by')} ${marker.name}</span>`"
                    >
                    </LPopup>
                </LMarker>
            </LLayerGroup>

            <LControl position="bottomleft" class="leaflet-control-attribution mouseposition">
                <b>{{ $t('common.latitude') }}</b
                >: {{ mouseLat }} | <b>{{ $t('common.longitude') }}</b
                >: {{ mouseLong }}
            </LControl>
            <LControl position="topleft">
                <div class="form-control flex flex-col gap-2">
                    <div>
                        <input
                            v-model="playerQuery"
                            class="w-full"
                            type="text"
                            name="searchPlayer"
                            :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
                        />
                    </div>
                    <div>
                        <input
                            v-model="dispatchQuery"
                            class="w-full"
                            type="text"
                            name="searchDispatch"
                            :placeholder="`${$t('common.dispatch', 1)} ${$t('common.filter')}`"
                        />
                    </div>
                    <div>
                        <Combobox as="div" class="w-full" v-model="selectedPostal" nullable>
                            <ComboboxInput
                                class="w-full"
                                @change="postalQuery = $event.target.value"
                                @click="loadPostals"
                                :display-value="(postal: any) => postal ? postal?.code : ''"
                                :placeholder="`${$t('common.postal')} ${$t('common.search')}`"
                            />
                            <ComboboxOptions class="z-10 w-full py-1 mt-1 overflow-auto bg-white">
                                <ComboboxOption
                                    v-for="postal in filteredPostals"
                                    :key="postal.code"
                                    :value="postal"
                                    v-slot="{ active }"
                                >
                                    <li
                                        :class="[
                                            'relative cursor-default select-none py-2 pl-8 pr-4',
                                            active ? 'bg-primary-500 text-white' : 'text-gray-600',
                                        ]"
                                    >
                                        {{ postal.code }}
                                    </li>
                                </ComboboxOption>
                            </ComboboxOptions>
                        </Combobox>
                    </div>
                </div>
            </LControl>
            <LControl position="bottomright">
                <div class="form-control flex flex-col gap-2">
                    <div class="p-2 bg-neutral border border-[#6b7280] flex flex-row justify-center">
                        <span class="text-lg mr-2 text-[#6f7683]">{{ $t('components.livemap.center_selected_marker') }}</span>
                        <input v-model="livemapCenterSelectedMarker" class="my-auto" name="livemapMarkerSize" type="checkbox" />
                    </div>
                    <div class="p-2 bg-neutral border border-[#6b7280] flex flex-row justify-center">
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
        </LMap>
    </div>
</template>

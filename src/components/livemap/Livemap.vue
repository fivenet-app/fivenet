<script lang="ts" setup>
import { onBeforeUnmount, ref } from 'vue';
import { ClientReadableStream, RpcError } from 'grpc-web';
import { StreamRequest, StreamResponse } from '@fivenet/gen/services/livemapper/livemap_pb';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import { LMap, LLayerGroup, LTileLayer, LControlLayers, LMarker, LPopup, LControl } from '@vue-leaflet/vue-leaflet';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import { ValueOf } from '~/utils/types';
import { DispatchMarker, UserMarker } from '@fivenet/gen/resources/livemap/livemap_pb';
import { Job } from '@fivenet/gen/resources/jobs/jobs_pb';
import { watchDebounced } from '@vueuse/core';
import { dispatchNotification } from '../partials/notification';
import { Combobox, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';

const { $grpc } = useNuxtApp();

const stream = ref<ClientReadableStream<StreamResponse> | null>(null);
const error = ref<RpcError | null>(null);

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

const markerJobs = ref<Job[]>([]);

const playerQuery = ref<string>('');
let playerMarkers: UserMarker[] = [];
const playerMarkersFiltered = ref<UserMarker[]>([]);

const dispatchQuery = ref<string>('');
let dispatchMarkers: DispatchMarker[] = [];
const dispatchMarkersFiltered = ref<DispatchMarker[]>([]);

async function applyPlayerQuery(): Promise<void> { playerMarkersFiltered.value = playerMarkers.filter(m => (m.getUser()?.getFirstname() + ' ' + m.getUser()?.getLastname()).includes(playerQuery.value)) }
async function applyDispatchQuery(): Promise<void> { dispatchMarkersFiltered.value = dispatchMarkers.filter(m => m.getPopup().includes(dispatchQuery.value) || m.getName().includes(dispatchQuery.value)) }

watchDebounced(playerQuery, async () => { applyPlayerQuery() }, { debounce: 750, maxWait: 2000 });
watchDebounced(dispatchQuery, async () => { applyDispatchQuery() }, { debounce: 750, maxWait: 2000 });

const mouseLat = ref<string>((0).toFixed(3));
const mouseLong = ref<string>((0).toFixed(3));

const currentHash = ref<string>('');
const isMoving = ref<boolean>(false);

let map: L.Map | undefined = undefined;

watch(currentHash, () => {
    window.location.replace(currentHash.value);
});

watchDebounced(isMoving, () => {
    if (!map || isMoving.value) return;

    const newHash = stringifyHash(map.getZoom(), map.getCenter().lat, map.getCenter().lng);
    if (currentHash.value !== newHash) currentHash.value = newHash;
}, { debounce: 1000, maxWait: 3000 });


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

function parseHash(hash: string): { latlng: L.LatLng, zoom: number } | undefined {
    if (hash.indexOf('#') === 0) hash = hash.substring(1);

    const args = hash.split('/');
    if (args.length !== 3) return;

    const zoom = parseInt(args[0], 10)
    const lat = parseFloat(args[1])
    const lng = parseFloat(args[2]);

    if (isNaN(zoom) || isNaN(lat) || isNaN(lng)) return;

    return {
        latlng: new L.LatLng(lat, lng),
        zoom,
    };
}

async function onMapReady($event: any): Promise<void> {
    map = $event as L.Map;

    const startingHash = window.location.hash;
    const startPos = parseHash(startingHash);
    if (startPos) $event.setView(startPos.latlng, startPos.zoom);

    $event.on('baselayerchange', async (event: L.LayersControlEvent) => { updateBackground(event.name) });

    $event.addEventListener('mousemove', async (event: L.LeafletMouseEvent) => {
        mouseLat.value = (Math.round(event.latlng.lat * 100000) / 100000).toFixed(3);
        mouseLong.value = (Math.round(event.latlng.lng * 100000) / 100000).toFixed(3);
    });

    $event.on('movestart', async () => { isMoving.value = true });

    $event.on('moveend', async () => { isMoving.value = false });

    startDataStream();
}

async function startDataStream(): Promise<void> {
    if (stream.value !== null) return;

    console.debug('Starting Data Stream');

    const request = new StreamRequest();

    stream.value = $grpc.getLivemapperClient().
        stream(request).
        on('error', async (err: RpcError) => {
            $grpc.handleRPCError(err);
            error.value = err;
            stopDataStream();
        }).
        on('data', async (resp) => {
            error.value = null;

            markerJobs.value = resp.getJobsList();

            playerMarkers = resp.getUsersList();
            dispatchMarkers = resp.getDispatchesList();

            applyPlayerQuery();
            applyDispatchQuery();
        }).
        on('end', async () => {
            console.debug('Data Stream Ended');
        });
}

async function stopDataStream(): Promise<void> {
    console.debug('Stopping Data Stream');
    if (stream.value !== null) {
        stream.value.cancel();
        stream.value = null;
    }
}

function getIcon(type: 'player' | 'dispatch', icon: string, iconColor: string): L.DivIcon {
    let html = ``;
    switch (type) {
        case 'player':
            {
                html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="${iconColor ? '#' + iconColor : 'currentColor'}" class="w-full h-full">
                    <path d="M8 16s6-5.686 6-10A6 6 0 0 0 2 6c0 4.314 6 10 6 10zm0-7a3 3 0 1 1 0-6 3 3 0 0 1 0 6z"/>
                </svg>`;
            }
            break;

        case 'dispatch':
            {
                html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="${iconColor ? '#' + iconColor : 'currentColor'}" class="w-full h-full">
                    <path d="M8 16a2 2 0 0 0 2-2H6a2 2 0 0 0 2 2zm.995-14.901a1 1 0 1 0-1.99 0A5.002 5.002 0 0 0 3 6c0 1.098-.5 6-2 7h14c-1.5-1-2-5.902-2-7 0-2.42-1.72-4.44-4.005-4.901z"/>
                </svg>`;
            }
            break;
    }

    return new L.DivIcon({
        html: '<div class="place-content-center">' + html + '</div>',
        iconSize: [48, 48],
        iconAnchor: [24, 24],
        popupAnchor: [-8, -24],
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
}

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
        postals.value.push(...(await response.json()) as Postal[]);
    } catch (_) {
        dispatchNotification({ title: 'Failed to load Postals map', content: '' });
        postalsLoaded = false;
    }
}

async function findPostal(): Promise<void> {
    if (postalQuery.value === "") {
        return;
    }

    let results = 0;
    filteredPostals.value.length = 0;
    filteredPostals.value = postals.value.filter(p => {
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

watch(selectedPostal, () => {
    if (!selectedPostal.value) {
        return;
    }

    map?.flyTo([selectedPostal.value.y, selectedPostal.value.x], 5, {
        duration: 0.850,
    });
});
watchDebounced(postalQuery, () => findPostal(), { debounce: 250, maxWait: 850 });
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
    transition: transform 1000ms ease;
}
</style>

<template>
    <div class="relative w-full h-full">
        <div v-if="error || stream === null" class="absolute inset-0 flex justify-center items-center"
            style="background-color: rgba(62, 60, 62, 0.5); z-index: 99999">
            <DataPendingBlock v-if="!error && stream === null" message="Starting Livemap data stream..." />
            <DataErrorBlock v-else-if="error" title="Failed to stream Livemap data!" :retry="() => { startDataStream() }" />
        </div>

        <LMap class="z-0" v-model:zoom="zoom" v-model:center="center" :crs="customCRS" :min-zoom="1" :max-zoom="6"
            :inertia="false" :style="{ backgroundColor }" @ready="onMapReady($event)" :use-global-leaflet="false">
            <LTileLayer url="/tiles/postal/{z}/{x}/{y}.png" layer-type="base" name="Postal" :no-wrap="true" :tms="true"
                :visible="true" :attribution="attribution" />
            <LTileLayer url="/tiles/atlas/{z}/{x}/{y}.png" layer-type="base" name="Atlas" :no-wrap="true" :tms="true"
                :visible="false" :attribution="attribution" />
            <LTileLayer url="/tiles/road/{z}/{x}/{y}.png" layer-type="base" name="Road" :no-wrap="true" :tms="true"
                :visible="false" :attribution="attribution" />
            <LTileLayer url="/tiles/satelite/{z}/{x}/{y}.png" layer-type="base" name="Satelite" :no-wrap="true" :tms="true"
                :visible="false" :attribution="attribution" />

            <LControlLayers />

            <LLayerGroup v-for="job in markerJobs" :key="job.getName()" :name="`Employees ${job.getLabel()}`"
                layer-type="overlay" :visible="true">
                <LMarker v-for="marker in playerMarkersFiltered.filter(p => p.getUser()?.getJob() === job.getName())"
                    :key="marker.getId()" :latLng="[marker.getY(), marker.getX()]" :name="marker.getName()"
                    :icon="getIcon('player', marker.getIcon(), marker.getIconColor()) as L.Icon">
                    <LPopup :options="{ closeButton: false }"
                        :content="`<span class='font-semibold'>Employee ${marker.getUser()?.getJobLabel()}</span><br><span class='italic'>[${marker.getUser()?.getJobGrade()}] ${marker.getUser()?.getJobGradeLabel()}</span><br>${marker.getUser()?.getFirstname()} ${marker.getUser()?.getLastname()}`">
                    </LPopup>
                </LMarker>
            </LLayerGroup>

            <LLayerGroup v-for="job in markerJobs" :key="job.getName()" :name="`Dispatches ${job.getLabel()}`"
                layer-type="overlay" :visible="true">
                <LMarker v-for="marker in dispatchMarkersFiltered.filter(m => m.getJob() === job.getName())"
                    :key="marker.getId()" :latLng="[marker.getY(), marker.getX()]" :name="marker.getName()"
                    :icon="getIcon('dispatch', marker.getIcon(), marker.getIconColor()) as L.Icon">
                    <LPopup :options="{ closeButton: false }"
                        :content="`<span class='font-semibold'>Dispatch ${marker.getJobLabel()}</span><br>${marker.getPopup()}<br><span class='italic'>Sent by ${marker.getName()}</span>`">
                    </LPopup>
                </LMarker>
            </LLayerGroup>

            <LControl position="bottomleft" class="leaflet-control-attribution mouseposition">
                <b>Latitude</b>: {{ mouseLat }} | <b>Longtiude</b>: {{ mouseLong }}
            </LControl>
            <LControl position="topleft">
                <div class="form-control">
                    <div class="relative flex items-center">
                        <input v-model="playerQuery" type="text" name="searchPlayer" id="searchPlayer"
                            placeholder="Employee Filter" />
                    </div>
                    <div class="relative flex items-center mt-2">
                        <input v-model="dispatchQuery" type="text" name="searchDispatch" id="searchDispatch"
                            placeholder="Dispatch Filter" />
                    </div>
                    <div class="relative flex items-center mt-2">
                        <Combobox as="div" v-model="selectedPostal" nullable>
                            <ComboboxInput @change="postalQuery = $event.target.value" @click="loadPostals"
                                :display-value="(postal: any) => postal ? postal?.code : ''" placeholder="Postal" />
                            <ComboboxOptions class="z-10 w-full py-1 mt-1 overflow-auto bg-white">
                                <ComboboxOption v-for="postal in filteredPostals" :key="postal.code" :value="postal"
                                    v-slot="{ active }">
                                    <li
                                        :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-primary-500 text-white' : 'text-gray-600']">
                                        {{ postal.code }}
                                    </li>
                                </ComboboxOption>
                            </ComboboxOptions>
                        </Combobox>
                    </div>
                </div>
            </LControl>
        </LMap>
    </div>
</template>

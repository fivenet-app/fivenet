<script lang="ts" setup>
import { onMounted, onBeforeUnmount, onUnmounted, ref, nextTick } from 'vue';
import { ClientReadableStream, RpcError } from 'grpc-web';
import { StreamRequest, StreamResponse } from '@fivenet/gen/services/livemapper/livemap_pb';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { customCRS, Livemap, MarkerType } from '../class/Livemap';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

const { $grpc } = useNuxtApp();

function round(num: number): number {
    return Math.round(num * 100) / 1000;
}

// Latitude and Longitiude popup on mouse over
let _latlng: HTMLDivElement;
const Position = L.Control.extend({
    _container: null,
    options: {
        position: 'bottomleft',
    },
    onAdd: function () {
        const latlng = L.DomUtil.create('div', 'leaflet-control-attribution mouseposition');
        _latlng = latlng;
        _latlng.innerHTML = '<b>Latitude</b>: 0.0 | <b>Longtiude</b>: 0.0';
        return latlng;
    },
    updateHTML: function (lat: number, lng: number) {
        _latlng.innerHTML = '<b>Latitude</b>: ' + round(lat) + ' | <b>Longtiude</b>: ' + round(lng);
    },
});
const position = new Position();

let map: Livemap | undefined = undefined;

const atlas = L.tileLayer('/tiles/atlas/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const postal = L.tileLayer('/tiles/postal/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const road = L.tileLayer('/tile/road/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const satelite = L.tileLayer('/tile/satelite/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});

let stream: ClientReadableStream<StreamResponse> | undefined = undefined;
const error = ref();

async function start(): Promise<void> {
    if (stream !== undefined) return;

    console.debug('Starting Livemap Data Stream');
    const request = new StreamRequest();

    stream = $grpc.getLivemapperClient().
        stream(request).
        on('error', async (err: RpcError) => {
            $grpc.handleRPCError(err);
            error.value = err;
            stop();
        }).
        on('data', async (resp) => {
            error.value = null;

            map?.parseMarkerlist(MarkerType.dispatch, resp.getDispatchesList());
            map?.parseMarkerlist(MarkerType.player, resp.getUsersList());
        }).
        on('end', function () {
            console.debug('Livemap Data Stream Ended');
        });
}

async function stop(): Promise<void> {
    console.debug('Stopping Livemap Data Stream');
    if (stream !== undefined) {
        stream.cancel();
        stream = undefined;
    }
}

const mapContainer = ref<HTMLDivElement | null>();

onMounted(() => {
    setTimeout(() => {
        if (!mapContainer.value) {
            return;
        }
        map = new Livemap(mapContainer.value, { layers: [postal], crs: customCRS });
        map.addHash();
        map.setView([0, 0], 2);

        const markersLayer = new L.LayerGroup().addTo(map as L.Map);
        L.control
            .layers({ Atlas: atlas, Road: road, Satelite: satelite, Postal: postal }, { Markers: markersLayer })
            .addTo(map as L.Map);
        postal.bringToFront();

        map.updateBackground('Postal');
        map.on('baselayerchange', (event: L.LayersControlEvent) => map?.updateBackground(event.name));

        map.addControl(position);
        map.addEventListener('mousemove', (event: L.LeafletMouseEvent) => {
            const lat = Math.round(event.latlng.lat * 100000) / 100000;
            const lng = Math.round(event.latlng.lng * 100000) / 100000;
            position.updateHTML(lat, lng);
        });

        start();
    }, 100);
});

onBeforeUnmount(() => {
    stop();
});

onUnmounted(() => {
    map = undefined;
});
</script>

<style>
.leaflet-div-icon {
    background: none;
    border: none;
}

.leaflet-div-icon svg path {
    stroke: #000000;
    stroke-width: 1.75px;
    stroke-linejoin: round;
}
</style>

<template>
    <div class="relative">
        <div id="map" ref="mapContainer" class="w-full z-0"></div>
        <div v-if="!stream || error" class="absolute inset-0 flex justify-center items-center z-10"
            style="background-color: rgba(62, 60, 62, 0.5)">
            <div v-if="error" class="rounded-md bg-red-50 p-4">
                <div class="flex">
                    <div class="flex-shrink-0">
                        <XCircleIcon class="h-5 w-5 text-red-400" aria-hidden="true" />
                    </div>
                    <div class="ml-3">
                        <h3 class="text-sm font-medium text-red-800">Failed to stream Livemap data</h3>
                        <div class="mt-2 text-sm text-red-700">
                            <p>
                                Please wait a few seconds and try again.
                            </p>
                        </div>
                        <div class="mt-4">
                            <div class="-mx-2 -my-1.5 flex">
                                <button type="button"
                                    class="rounded-md bg-red-50 px-2 py-1.5 text-sm font-medium text-red-800 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-red-600 focus:ring-offset-2 focus:ring-offset-red-50"
                                    @click="start()">
                                    Retry
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

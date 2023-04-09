<script lang="ts" setup>
import { onMounted, onBeforeUnmount, onUnmounted, ref, nextTick } from 'vue';
import { ClientReadableStream, RpcError } from 'grpc-web';
import { StreamRequest, StreamResponse } from '@fivenet/gen/services/livemapper/livemap_pb';
import L from 'leaflet';
import { LMap, LLayerGroup, LTileLayer, LControlLayers, LMarker, LPopup, LControl } from '@vue-leaflet/vue-leaflet';
import 'leaflet/dist/leaflet.css';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import { ValueOf } from '../../utils/types';
import { DispatchMarker, UserMarker } from '@fivenet/gen/resources/livemap/livemap_pb';
import { Job } from '@fivenet/gen/resources/jobs/jobs_pb';
import { watchDebounced } from '@vueuse/core';

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

const zoom: number = 2;
const center: L.PointExpression = [0, 0];
const attribution = '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>';

const markerJobs = ref<Job[]>([]);
const playerMarkers = ref<UserMarker[]>([]);
const dispatchMarkers = ref<DispatchMarker[]>([]);

const mouseLat = ref<string>((0).toFixed(3));
const mouseLong = ref<string>((0).toFixed(3));

const currentHash = ref<string>('');
const isMoving = ref<boolean>(false);

let map: L.Map | undefined = undefined;

watch(currentHash, () => {
    window.location.replace(currentHash.value);
})

watchDebounced(isMoving, () => {
    if (isMoving.value || !map) return;

    const newHash = stringifyHash(map.getZoom(), map.getCenter().lat, map.getCenter().lng);
    if (currentHash.value !== newHash) currentHash.value = newHash;
}, { debounce: 1000, maxWait: 3000 })


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
    const precision = Math.max(0, Math.ceil(Math.log(zoom) / Math.LN2));

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

    $event.on('movestart', async () => { isMoving.value = true })

    $event.on('moveend', async () => { isMoving.value = false })

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

            dispatchMarkers.value = resp.getDispatchesList();
            playerMarkers.value = resp.getUsersList();
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
    let html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-full h-full mx-auto">
          <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm11.378-3.917c-.89-.777-2.366-.777-3.255 0a.75.75 0 01-.988-1.129c1.454-1.272 3.776-1.272 5.23 0 1.513 1.324 1.513 3.518 0 4.842a3.75 3.75 0 01-.837.552c-.676.328-1.028.774-1.028 1.152v.75a.75.75 0 01-1.5 0v-.75c0-1.279 1.06-2.107 1.875-2.502.182-.088.351-.199.503-.331.83-.727.83-1.857 0-2.584zM12 18a.75.75 0 100-1.5.75.75 0 000 1.5z" clip-rule="evenodd" />
        </svg>`;
    switch (type) {
        case 'player':
            {
                html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="${iconColor ? '#' + iconColor : 'currentColor'
                    }" class="w-full h-full">
                  <path fill-rule="evenodd" d="M11.54 22.351l.07.04.028.016a.76.76 0 00.723 0l.028-.015.071-.041a16.975 16.975 0 001.144-.742 19.58 19.58 0 002.683-2.282c1.944-1.99 3.963-4.98 3.963-8.827a8.25 8.25 0 00-16.5 0c0 3.846 2.02 6.837 3.963 8.827a19.58 19.58 0 002.682 2.282 16.975 16.975 0 001.145.742zM12 13.5a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
                </svg>`;
            }
            break;

        case 'dispatch':
            {
                html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="${iconColor ? '#' + iconColor : 'currentColor'
                    }" class="w-full h-full">
                  <path fill-rule="evenodd" d="M5.25 9a6.75 6.75 0 0113.5 0v.75c0 2.123.8 4.057 2.118 5.52a.75.75 0 01-.297 1.206c-1.544.57-3.16.99-4.831 1.243a3.75 3.75 0 11-7.48 0 24.585 24.585 0 01-4.831-1.244.75.75 0 01-.298-1.205A8.217 8.217 0 005.25 9.75V9zm4.502 8.9a2.25 2.25 0 104.496 0 25.057 25.057 0 01-4.496 0z" clip-rule="evenodd" />
                </svg>`;
            }
            break;
    }

    return new L.DivIcon({
        html: '<div class="place-content-center">' + html + '</div>',
        iconSize: [36, 36],
        iconAnchor: [18, 18],
        popupAnchor: [0, -8],
    });
}

onBeforeUnmount(() => {
    stopDataStream();
});

// onMounted(() => {
//     setTimeout(async () => {
//         if (!mapContainer.value) {
//             return;
//         }
//         map = new Livemap(mapContainer.value, { layers: [postal], crs: customCRS });
//         map.addHash();
//         map.setView([0, 0], 2);

//         await map.addLayerGroup('Players');
//         await map.addLayerGroup('Dispatches');

//         await map.addControlLayer({ Atlas: atlas, Road: road, Satelite: satelite, Postal: postal });

//         postal.bringToFront();

//         map.updateBackground('Postal');
//         map.on('baselayerchange', (event: L.LayersControlEvent) => map?.updateBackground(event.name));

//         map.addControl(position);
//         map.addEventListener('mousemove', (event: L.LeafletMouseEvent) => {
//             const lat = Math.round(event.latlng.lat * 100000) / 100000;
//             const lng = Math.round(event.latlng.lng * 100000) / 100000;
//             position.updateHTML(lat, lng);
//         });

//         start();
//     }, 100);
// });
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

.leaflet-marker-icon {
    transition: transform 1000ms ease;
}
</style>

<template>
    <div class="w-full relative h-full">
        <div v-if="error || stream === null" class="absolute inset-0 flex justify-center items-center"
            style="background-color: rgba(62, 60, 62, 0.5); z-index: 99999">
            <DataPendingBlock v-if="!error && stream === null" message="Starting Livemap data stream..." />
            <DataErrorBlock v-else-if="error" title="Failed to stream Livemap data!" :retry="() => { startDataStream() }" />
        </div>

        <LMap ref="mapElement" v-model:zoom="zoom" v-model:center="center" :crs="customCRS" :min-zoom="0" :max-zoom="6"
            :inertia="false" :style="{ backgroundColor }" @ready="onMapReady($event)">
            <LTileLayer url="/tiles/postal/{z}/{x}/{y}.png" layer-type="base" name="Postal" :no-wrap="true" :tms="true"
                :visible="true" :attribution="attribution" />
            <LTileLayer url="/tiles/atlas/{z}/{x}/{y}.png" layer-type="base" name="Atlas" :no-wrap="true" :tms="true"
                :visible="false" :attribution="attribution" />
            <LTileLayer url="/tiles/road/{z}/{x}/{y}.png" layer-type="base" name="Road" :no-wrap="true" :tms="true"
                :visible="false" :attribution="attribution" />
            <LTileLayer url="/tiles/satelite/{z}/{x}/{y}.png" layer-type="base" name="Satelite" :no-wrap="true" :tms="true"
                :visible="false" :attribution="attribution" />

            <LControlLayers />

            <LLayerGroup v-for="job in markerJobs" :key="job.getName()" :name="`Players ${job.getLabel()}`"
                layer-type="overlay" :visible="true">
                <LMarker v-for="marker in playerMarkers.filter(p => p.getUser()?.getJob() === job.getName())"
                    :key="marker.getId()" :latLng="[marker.getY(), marker.getX()]" :name="marker.getName()"
                    :icon="getIcon('player', marker.getIcon(), marker.getIconColor()) as L.Icon">
                    <LPopup :options="{ closeButton: false }"
                        :content="`${marker.getUser()?.getFirstname()}, ${marker.getUser()?.getLastname()} (Job: ${marker.getUser()?.getJobLabel()})`">
                    </LPopup>
                </LMarker>
            </LLayerGroup>

            <LLayerGroup v-for="job in markerJobs" :key="job.getName()" :name="`Dispatches ${job.getLabel()}`"
                layer-type="overlay" :visible="true">
                <LMarker v-for="marker in dispatchMarkers.filter(m => m.getJob() === job.getName())" :key="marker.getId()"
                    :latLng="[marker.getY(), marker.getX()]" :name="marker.getName()"
                    :icon="getIcon('dispatch', marker.getIcon(), marker.getIconColor()) as L.Icon">
                    <LPopup :options="{ closeButton: false }"
                        :content="`Dispatch: ${marker.getPopup()}<br>Sent by: ${marker.getName()} (Job: ${marker.getJobLabel()})`">
                    </LPopup>
                </LMarker>
            </LLayerGroup>

            <LControl position="bottomleft" class="leaflet-control-attribution mouseposition">
                <b>Latitude</b>: {{ mouseLat }} | <b>Longtiude</b>: {{ mouseLong }}
            </LControl>
        </LMap>
    </div>
</template>

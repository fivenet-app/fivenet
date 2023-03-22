<script lang="ts" setup>
import { onMounted, onBeforeUnmount, onUnmounted } from 'vue';
import { getLivemapperClient } from '../grpc/grpc';
import { ClientReadableStream, RpcError } from 'grpc-web';
import { StreamRequest, StreamResponse } from '@arpanet/gen/services/livemapper/livemap_pb';
// Leaflet and Livemap custom parts
import { customCRS, Livemap, MarkerType } from '../class/Livemap';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { handleGRPCError } from '../grpc/interceptors';

// Latitude and Longitiude popup on mouse over
let _latlng: HTMLDivElement;
const Position = L.Control.extend({
    _container: null,
    options: {
        position: 'bottomleft',
    },
    onAdd: function () {
        const latlng = L.DomUtil.create('div', 'mouseposition');
        _latlng = latlng;
        return latlng;
    },
    updateHTML: function (lat: number, lng: number) {
        _latlng.innerHTML = 'Latitude: ' + lat + '   Longitiude: ' + lng;
    },
});
const position = new Position();

let stream = null as null | ClientReadableStream<StreamResponse>;
let map = {} as undefined | Livemap;

const atlas = L.tileLayer('tiles/atlas/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const postal = L.tileLayer('tiles/postal/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const road = L.tileLayer('tiles/road/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const satelite = L.tileLayer('tiles/satelite/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});

function start() {
    console.log('starting livemap data stream');
    const request = new StreamRequest();

    stream = getLivemapperClient()
        .stream(request)
        .on('error', (err: RpcError) => {
            handleGRPCError(err);
        })
        .on('data', function (resp) {
            map?.parseMarkerlist(MarkerType.dispatch, resp.getDispatchesList());
            map?.parseMarkerlist(MarkerType.player, resp.getUsersList());
        })
        .on('end', function () {
            console.log('livemap data stream ended');
        });
}

function stop() {
    console.log('stopping livemap data stream');
    if (stream) {
        stream.cancel();
        stream = null;
    }
}

onMounted(() => {
    map = new Livemap('map', { layers: [postal], crs: customCRS });
    map.addHash();
    map.setView([0, 0], 2);

    const markersLayer = new L.LayerGroup().addTo(map as L.Map);
    L.control
        .layers({ Satelite: satelite, Atlas: atlas, Road: road, Postal: postal }, { Markers: markersLayer })
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
});

onBeforeUnmount(() => {
    stop();
});

onUnmounted(() => {
    map = undefined;
});
</script>

<style scoped>
#map {
    height: 94vh;
}
</style>

<template>
    <div class="w-full z-0" id="map"></div>
</template>

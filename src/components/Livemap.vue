<script lang="ts">
import { defineComponent } from 'vue';
import authInterceptor from '../grpcauth';
import { ClientReadableStream, RpcError } from 'grpc-web';
import { LivemapServiceClient } from '@arpanet/gen/livemap/LivemapServiceClientPb';
import { Marker, StreamRequest, ServerStreamResponse } from '@arpanet/gen/livemap/livemap_pb';
// Leaflet and Livemap custom parts
import { customCRS, Livemap } from '../class/Livemap';
import { Hash } from '../class/Hash';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

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

export default defineComponent({
    data() {
        return {
            client: new LivemapServiceClient('https://localhost:8181', null, {
                unaryInterceptors: [authInterceptor],
                streamInterceptors: [authInterceptor],
            }),
            stream: null as null | ClientReadableStream<ServerStreamResponse>,
            map: {} as Livemap,
            hash: {} as Hash,
            zoom: 1,
            usersList: [] as Array<Marker>,
            dispatchesList: [] as Array<Marker>,
        };
    },
    beforeUnmount: function () {
        this.stop();
    },
    setup: function () {
        return {
            customCRS,
            position,
        };
    },
    mounted() {
        this.map = new Livemap('map', { layers: [postal], crs: customCRS });
        this.map.addHash();
        this.map.setView([0, 0], 2);

        const markersLayer = new L.LayerGroup().addTo(this.map as L.Map);
        L.control
            .layers({ Satelite: satelite, Atlas: atlas, Road: road, Postal: postal }, { Markers: markersLayer })
            .addTo(this.map as L.Map);
        postal.bringToFront();

        this.updateBackground('Postal');
        this.map.on('baselayerchange', (context) => this.updateBackground(context.name));

        this.map.addControl(position);
        this.map.addEventListener('mousemove', (event: L.LeafletMouseEvent) => {
            const lat = Math.round(event.latlng.lat * 100000) / 100000;
            const lng = Math.round(event.latlng.lng * 100000) / 100000;
            position.updateHTML(lat, lng);
        });

        this.start();
    },
    methods: {
        updateBackground(layer: string): void {
            switch (layer) {
                case 'Atlas':
                    this.map.getContainer().style.backgroundColor = '#0fa8d2';
                    return;
                case 'Satelite':
                    this.map.getContainer().style.backgroundColor = '#143d6b';
                    return;
                case 'Road':
                    this.map.getContainer().style.backgroundColor = '#1862ad';
                    return;
                case 'Postal':
                    this.map.getContainer().style.backgroundColor = '#63a7ce';
                    return;
            }
        },
        start: function () {
            console.log('starting livemap data stream');

            let outer = this;
            const request = new StreamRequest();
            this.stream = this.client
                .stream(request)
                .on('data', function (response) {
                    outer.usersList = response.getUsersList();
                    // TODO Marker management
                })
                .on('error', (err: RpcError) => {
                    authInterceptor.handleError(err, this.$route);
                })
                .on('end', function () {
                    console.log('livemap data stream ended');
                });
        },
        stop: function () {
            console.log('stopping livemap data stream');
            if (this.stream) {
                this.stream.cancel();
                this.stream = null;
            }
        },
    },
});

const atlas = L.tileLayer('tiles/atlas/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const road = L.tileLayer('tiles/road/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const satelite = L.tileLayer('tiles/satelite/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
const postal = L.tileLayer('tiles/postal/{z}/{x}/{y}.png', {
    attribution:
        '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
    minZoom: 1,
    maxZoom: 6,
    noWrap: false,
    tms: true,
});
</script>

<template>
    <div class="w-screen" id="map" style="height: 95vh"></div>
</template>

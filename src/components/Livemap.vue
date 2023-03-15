<script lang="ts">
import { defineComponent } from 'vue';
import { getLivemapperClient, handleGRPCError } from '../grpc';
import { ClientReadableStream, RpcError } from 'grpc-web';
import { StreamRequest, ServerStreamResponse } from '@arpanet/gen/services/livemapper/livemap_pb';
import { DispatchMarker, UserMarker } from '@arpanet/gen/resources/livemap/livemap_pb';
// Leaflet and Livemap custom parts
import { customCRS, Livemap, MarkerType } from '../class/Livemap';
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
            stream: null as null | ClientReadableStream<ServerStreamResponse>,
            map: {} as undefined | Livemap,
            hash: {} as Hash,
            zoom: 1,
            dispatchesList: [] as Array<DispatchMarker>,
            usersList: [] as Array<UserMarker>,
        };
    },
    setup: function () {
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

        return {
            customCRS,
            position,
            atlas,
            postal,
            road,
            satelite,
        };
    },
    mounted() {
        this.map = new Livemap('map', { layers: [this.postal], crs: customCRS });
        this.map.addHash();
        this.map.setView([0, 0], 2);

        const markersLayer = new L.LayerGroup().addTo(this.map as L.Map);
        L.control
            .layers({ Satelite: this.satelite, Atlas: this.atlas, Road: this.road, Postal: this.postal }, { Markers: markersLayer })
            .addTo(this.map as L.Map);
        this.postal.bringToFront();

        this.map.updateBackground('Postal');
        this.map.on('baselayerchange', (context) => this.map?.updateBackground(context.name));

        this.map.addControl(position);
        this.map.addEventListener('mousemove', (event: L.LeafletMouseEvent) => {
            const lat = Math.round(event.latlng.lat * 100000) / 100000;
            const lng = Math.round(event.latlng.lng * 100000) / 100000;
            position.updateHTML(lat, lng);
        });

        this.start();
    },
    beforeUnmount: function () {
        this.stop();
    },
    unmounted() {
        this.map = undefined;
    },
    methods: {
        start: function () {
            console.log('starting livemap data stream');

            let outer = this;
            const request = new StreamRequest();

            this.stream = getLivemapperClient()
                .stream(request)
                .on('data', function (resp) {
                    outer.usersList = resp.getUsersList();
                    outer.map?.parseMarkerlist(MarkerType.player, outer.usersList);
                })
                .on('error', (err: RpcError) => {
                    handleGRPCError(err, this.$route);
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
        updateBackground: function (layer: string) {
            this.map?.updateBackground(layer);
        },
    },
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

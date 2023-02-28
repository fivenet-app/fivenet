<script lang="ts">
import { defineComponent } from 'vue';
import { mapState } from 'vuex';
import authInterceptor from '../grpcauth';
import * as grpcWeb from 'grpc-web';
import { LivemapServiceClient } from '@arpanet/gen/livemap/LivemapServiceClientPb';
import { Marker, StreamRequest, ServerStreamResponse } from '@arpanet/gen/livemap/livemap_pb';
// Leaflet and Livemap custom parts
import { LMap, LTileLayer, LMarker, LControlLayers, LLayerGroup, LPopup, LControlScale } from "@vue-leaflet/vue-leaflet";
import { customCRS } from '../livemap/CRS';
import { Hash } from '../livemap/Hash';
import L from 'leaflet';

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
    components: {
        LMap,
        LTileLayer,
        LMarker,
        LControlLayers,
        LLayerGroup,
        LPopup,
        LControlScale,
    },
    computed: {
        ...mapState({
            accessToken: 'accessToken',
        }),
    },
    data: function () {
        return {
            map: {} as LMap,
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
        }
    },
    mounted() {
        this.map = this.$refs.map as LMap;
    },
    methods: {
        onLeafletReady: function () {
            this.map.leafletObject.on('baselayerchange', (context: L.LayersControlEvent) => this.updateBackground(context.name));
            this.updateBackground('Postal');

            // Register Coordinates Display
            // TODO This should probably be a "sub"-component that dynamically listens to the lat/lng changes
            this.map.leafletObject.addControl(position);
            this.map.leafletObject.addEventListener('mousemove', (event: L.LeafletMouseEvent) => {
                const lat = Math.round(event.latlng.lat * 100000) / 100000;
                const lng = Math.round(event.latlng.lng * 100000) / 100000;
                position.updateHTML(lat, lng);
            });

            // Map Position Hash
            //this.hash = new Hash(this.map.leafletObject, this.map.leafletObject.getContainer());

            this.start();
        },
        updateBackground(layer: string): void {
            switch (layer) {
                case 'Atlas':
                    this.map.leafletObject.getContainer().style.backgroundColor = '#0fa8d2';
                    return;
                case 'Satelite':
                    this.map.leafletObject.getContainer().style.backgroundColor = '#143d6b';
                    return;
                case 'Road':
                    this.map.leafletObject.getContainer().style.backgroundColor = '#1862ad';
                    return;
                case 'Postal':
                    this.map.leafletObject.getContainer().style.backgroundColor = '#8cb6ce';
                    return;
            }
        },
        start: function () {
            console.log("starting livemap data stream");

            let outer = this;
            const request = new StreamRequest();
            stream = client.stream(request);
            stream.on('data', function (response) {
                outer.usersList = response.getUsersList();
            });
            stream.on('end', function () {
                console.log('livemap data stream ended');
            });
        },
        stop: function () {
            console.log("stopping livemap data stream");
            stream.cancel();
        },
    }
});

const client = new LivemapServiceClient('https://localhost:8181', null, {
    unaryInterceptors: [authInterceptor],
    streamInterceptors: [authInterceptor],
});
let stream: grpcWeb.ClientReadableStream<ServerStreamResponse>;
</script>

<template>
    <div class="w-screen" style="height: 95vh">
        <l-map ref="map" v-model:zoom="zoom" :center="[0, 0]" :crs="customCRS" :registerControl="position"
            @ready="onLeafletReady">
            <l-tile-layer url="tiles/atlas/{z}/{x}/{y}.png" layer-type="base" :tms=true :no-wrap=false name="Atlas"
                attribution="<a href='http://www.rockstargames.com/V/'>Grand Theft Auto V</a>" :min-zoom="0"
                :max-zoom="6"></l-tile-layer>
            <l-tile-layer url="tiles/road/{z}/{x}/{y}.png" layer-type="base" :tms=true :no-wrap=false name="Road"
                attribution="<a href='http://www.rockstargames.com/V/'>Grand Theft Auto V</a>" :min-zoom="0"
                :max-zoom="6"></l-tile-layer>
            <l-tile-layer url="tiles/satelite/{z}/{x}/{y}.png" layer-type="base" :tms=true :no-wrap=false name="Satelite"
                attribution="<a href='http://www.rockstargames.com/V/'>Grand Theft Auto V</a>" :min-zoom="0"
                :max-zoom="6"></l-tile-layer>
            <l-tile-layer url="tiles/postal/{z}/{x}/{y}.png" layer-type="base" :tms=true :no-wrap=false name="Postals"
                attribution="<a href='http://www.rockstargames.com/V/'>Grand Theft Auto V</a>" :min-zoom="0"
                :max-zoom="6"></l-tile-layer>
            <l-layer-group name="Markers">
                <l-marker v-for="marker in usersList" :lat-lng="[marker.getX(), marker.getY()]" :name="marker.getName()"
                    :draggable=false>
                    <l-popup>{{ marker.getPopup() ? marker.getPopup() : marker.getName() }}</l-popup>
                </l-marker>
            </l-layer-group>
            <l-control-layers />
            <l-control-scale />
        </l-map>
    </div>
</template>

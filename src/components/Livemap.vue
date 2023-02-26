<script lang="ts">
import { defineComponent } from 'vue';
import * as grpcWeb from 'grpc-web';
import { LivemapServiceClient } from '@arpanet/gen/livemap/LivemapServiceClientPb';
import { Marker, StreamRequest, ServerStreamResponse } from '@arpanet/gen/livemap/livemap_pb';
import { LMap, LTileLayer, LMarker, LControlLayers, LLayerGroup, LPopup } from "@vue-leaflet/vue-leaflet";
import { customCRS } from '../livemap/CRS';
import { Hash } from '../livemap/Hash';
import L from 'leaflet';

const service = new LivemapServiceClient('https://localhost:8181', null, null);
let stream: grpcWeb.ClientReadableStream<ServerStreamResponse>;

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

            // Register Coordinates Display
            // TODO This should probably be a "sub"-component that dynamically listens to the lat/lng changes
            this.map.leafletObject.addControl(position);
            this.map.leafletObject.addEventListener('mousemove', (event: L.LeafletMouseEvent) => {
                const lat = Math.round(event.latlng.lat * 100000) / 100000;
                const lng = Math.round(event.latlng.lng * 100000) / 100000;
                position.updateHTML(lat, lng);
            });

            // Map Position Hash
            this.hash = new Hash(this.map.leafletObject, this.map.leafletObject.getContainer());

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
            }
        },
        start: function () {
            console.log("starting livemap data stream");

            let outer = this;
            const request = new StreamRequest();
            stream = service.stream(request);
            stream.on('data', function (response) {
                console.log("livemap data received");
                outer.usersList = response.getUsersList();
            });
            stream.on('end', function () {
                console.log('livemap data stream ended');
            });

            console.log("started livemap data stream");
        },
        stop: function () {
            console.log("stopping livemap data stream");
            stream.cancel();
            console.log("stopped livemap data stream");
        },
    }
});
</script>

<template>
    <div style="height:1000px; width:1920px">
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
            <l-layer-group name="Markers">
                <l-marker v-for="marker in usersList" :lat-lng="[marker.getX(), marker.getY()]" :name="marker.getName()"
                    :draggable=false>
                    <l-popup>{{ marker.getPopup() ? marker.getPopup() : marker.getName() }}</l-popup>
                </l-marker>
            </l-layer-group>
            <l-control-layers />
        </l-map>
    </div>
</template>

<style>
html,
body {
    height: 100%;
    min-width: 100%;
}

#app {
    min-height: 100%;
    min-width: 100%;
}

body {
    margin: 0;
    font-family: Helvetica, Arial, sans-serif;
    font-size: 12px;
    overflow: hidden;
}
</style>

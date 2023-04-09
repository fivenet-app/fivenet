import { IMarker } from '@fivenet/gen/resources/livemap/livemap';
import { DispatchMarker, UserMarker } from '@fivenet/gen/resources/livemap/livemap_pb';
import L from 'leaflet';

import { Hash } from './Hash';

export enum MarkerType {
    'player',
    'dispatch',
}

export class Livemap extends L.Map {
    public hash: Hash | undefined;
    public hasLoaded: boolean = false;

    public layerGroups: Map<string, L.LayerGroup> = new Map();
    public controlLayer: L.Control.Layers | undefined = undefined;

    public markers: Map<number, L.Marker> = new Map();
    public markersQuery: string = '';
    public markersFiltered: typeof this.markers = new Map();

    public popups: Map<number, L.Popup> = new Map();
    private prevMarkerLists: Map<MarkerType, Array<IMarker.AsObject>> = new Map();
    private defaultIcon: L.Icon;

    private element: HTMLElement;

    constructor(element: HTMLElement, options: L.MapOptions) {
        super(element, options);
        this.element = element;

        this.defaultIcon = new L.Icon({
            iconUrl: import.meta.env.BASE_URL + 'images/livemap/markers/user-default.svg',
            iconSize: [36, 36],
            iconAnchor: [18, 18],
            popupAnchor: [0, -8],
        });

        this.on('load', () => (this.hasLoaded = true));
        this.on('baselayerchange', (context) => this.updateBackground(context.name));
    }

    public async addHash(): Promise<void> {
        this.hash = new Hash(this, this.element);
    }

    public async removeHash(): Promise<boolean> {
        if (!this.hash) return false;

        this.hash.remove();
        return true;
    }

    public async updateBackground(layer: string): Promise<void> {
        switch (layer) {
            case 'Atlas':
                this.element.style.backgroundColor = '#0fa8d2';
                return;
            case 'Satelite':
                this.element.style.backgroundColor = '#143d6b';
                return;
            case 'Road':
                this.element.style.backgroundColor = '#1862ad';
                return;
            case 'Postal':
                this.element.style.backgroundColor = '#74aace';
                return;
        }
    }

    public async addLayerGroup(name: string): Promise<L.LayerGroup> {
        const layer = new L.LayerGroup();
        layer.addTo(this);

        this.layerGroups.set(name, layer);

        return layer;
    }

    public async getLayerGroup(id: string): Promise<L.LayerGroup | undefined> {
        return this.layerGroups.get(id);
    }

    public async addControlLayer(baseLayers: L.Control.LayersObject): Promise<L.Control.Layers | undefined> {
        const overlayLayers = Object.fromEntries(this.layerGroups.entries());
        this.controlLayer = L.control.layers(baseLayers, overlayLayers).addTo(this);

        return this.controlLayer;
    }

    public async getControlLayer(id: string): Promise<L.Control.Layers | undefined> {
        return this.controlLayer;
    }

    public async addMarker(
        id: number,
        latitude: number,
        longitude: number,
        popupContent: string,
        options: L.MarkerOptions,
        layerName: string
    ): Promise<L.Marker> {
        let marker = this.markers.get(id);
        const layer = this.layerGroups.get(layerName) ?? this;

        if (marker) {
            marker.setLatLng(L.latLng(latitude, longitude));
            if (options?.icon) marker.setIcon(options.icon);
            if (options?.opacity) marker.setOpacity(options.opacity);
        } else {
            options.icon = options?.icon ? options.icon : this.defaultIcon;
            options.icon.options.shadowSize = [0, 0];

            marker = new L.Marker(L.latLng(latitude, longitude), options).addTo(layer);

            if (popupContent) {
                const popup = L.popup({ content: popupContent, closeButton: false });
                this.popups.set(id, popup);
                marker.bindPopup(popup);
            }

            this.markers.set(id, marker);
        }

        return marker;
    }

    public async removeMarker(id: number): Promise<boolean> {
        const marker = this.markers.get(id);
        if (!marker) return false;

        marker.remove();
        return this.markers.delete(id);
    }

    public async parseMarkerlist(type: MarkerType, list: IMarker[]): Promise<void> {
        let options: L.MarkerOptions = {};
        let layer: string = '';

        switch (type) {
            case MarkerType.player:
                {
                    layer = 'Players';
                    options = {};
                }
                break;

            case MarkerType.dispatch:
                {
                    layer = 'Dispatches';
                    options = {};
                }
                break;
        }

        const previousList = this.prevMarkerLists.get(type);
        if (previousList) {
            const markersToRemove = previousList.filter((entry) => !list.find((e) => e.getId() === entry.id));
            markersToRemove.forEach((marker) => {
                this.removeMarker(marker.id);
            });
        }

        list.forEach(async (marker) => {
            if (marker.getIcon() || marker.getIconColor())
                options.icon = await this.getIcon(type, marker.getIcon(), marker.getIconColor());
            let popupContent = '';
            if (type === MarkerType.player) {
                const userMarker = marker as UserMarker;
                popupContent += `${userMarker.getUser()?.getFirstname()}, ${userMarker
                    .getUser()
                    ?.getLastname()} (Job: ${userMarker.getUser()?.getJobLabel()})`;
            } else if (type === MarkerType.dispatch) {
                const dispatchMarker = marker as DispatchMarker;
                popupContent += `Dispatch: ${dispatchMarker.getPopup()}<br>Sent by: ${dispatchMarker.getName()} (Job: ${dispatchMarker.getJobLabel()})`;
            }
            await this.addMarker(marker.getId(), marker.getY(), marker.getX(), popupContent, options, layer);
        });

        this.prevMarkerLists.set(
            type,
            list.map((e) => e.toObject())
        );
    }

    public async getIcon(type: MarkerType, icon: string, iconColor: string): Promise<L.DivIcon> {
        let html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-full h-full mx-auto">
          <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm11.378-3.917c-.89-.777-2.366-.777-3.255 0a.75.75 0 01-.988-1.129c1.454-1.272 3.776-1.272 5.23 0 1.513 1.324 1.513 3.518 0 4.842a3.75 3.75 0 01-.837.552c-.676.328-1.028.774-1.028 1.152v.75a.75.75 0 01-1.5 0v-.75c0-1.279 1.06-2.107 1.875-2.502.182-.088.351-.199.503-.331.83-.727.83-1.857 0-2.584zM12 18a.75.75 0 100-1.5.75.75 0 000 1.5z" clip-rule="evenodd" />
        </svg>`;
        switch (type) {
            case MarkerType.player:
                {
                    html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="${
                        iconColor ? '#' + iconColor : 'currentColor'
                    }" class="w-full h-full">
                  <path fill-rule="evenodd" d="M11.54 22.351l.07.04.028.016a.76.76 0 00.723 0l.028-.015.071-.041a16.975 16.975 0 001.144-.742 19.58 19.58 0 002.683-2.282c1.944-1.99 3.963-4.98 3.963-8.827a8.25 8.25 0 00-16.5 0c0 3.846 2.02 6.837 3.963 8.827a19.58 19.58 0 002.682 2.282 16.975 16.975 0 001.145.742zM12 13.5a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
                </svg>`;
                }
                break;

            case MarkerType.dispatch:
                {
                    html = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="${
                        iconColor ? '#' + iconColor : 'currentColor'
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
}

const centerX = 117.3;
const centerY = 172.8;
const scaleX = 0.02072;
const scaleY = 0.0205;

export const customCRS = L.extend({}, L.CRS.Simple, {
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

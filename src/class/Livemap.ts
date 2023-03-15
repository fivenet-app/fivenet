import { UserMarker } from '@arpanet/gen/resources/livemap/livemap_pb';
import L from 'leaflet';

import { AnimatedMarker } from './AnimatedMarker';
import { Hash } from './Hash';

export enum MarkerType {
    'player',
    'dispatch',
}

export class Livemap extends L.Map {
    public hash: Hash | undefined;
    public hasLoaded: boolean = false;

    public markers: Map<number, AnimatedMarker> = new Map();
    public popups: Map<number, L.Popup> = new Map();
    private prevMarkerLists: Map<MarkerType, Array<UserMarker.AsObject>> = new Map();

    private element: HTMLElement;

    constructor(element: string | HTMLElement, options?: L.MapOptions | undefined) {
        super(element, options);
        this.element = typeof element === 'string' ? (document.getElementById(element) as HTMLElement) : element;

        this.on('load', () => (this.hasLoaded = true));
        this.on('baselayerchange', (context) => this.updateBackground(context.name));
    }

    public addHash(): void {
        this.hash = new Hash(this, this.element);
    }

    public removeHash(): boolean {
        if (!this.hash) return false;

        this.hash.remove();
        return true;
    }

    public updateBackground(layer: string): void {
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

    public addMarker(id: number, latitude: number, longitude: number, content: string, options: L.MarkerOptions): void {
        const marker = this.markers.get(id);

        if (marker) {
            marker.moveTo(L.latLng(latitude, longitude));
            if (options?.icon) marker.setIcon(options.icon);
            if (options?.opacity) marker.setOpacity(options.opacity);
        } else {
            options.icon = options?.icon ? options.icon : new L.Icon.Default();
            options.icon.options.shadowSize = [0, 0];

            const marker = new AnimatedMarker(L.latLng(latitude, longitude), options).addTo(this);

            if (content) {
                const popup = L.popup({ content, closeButton: false });
                this.popups.set(id, popup);
                marker.bindPopup(popup);
            }

            this.markers.set(id, marker);
        }
    }

    public removeMarker(id: number): boolean {
        const marker = this.markers.get(id);
        if (!marker) return false;

        marker.remove();
        return this.markers.delete(id);
    }

    public parseMarkerlist(type: MarkerType, list: Array<UserMarker>): void {
        let options: L.MarkerOptions = {};
        switch (type) {
            case MarkerType.player: {
                options = {};
            }

            case MarkerType.dispatch: {
                options = {};
            }
        }

        const previousList = this.prevMarkerLists.get(type);
        if (previousList) {
            const markersToRemove = previousList.filter((entry) => !list.find((e) => e.getUserId() === entry.userId));
            markersToRemove.forEach((marker) => {
                this.removeMarker(marker.userId);
            });
        }

        list.forEach((marker) => {
            this.addMarker(marker.getUserId(), marker.getY(), marker.getX(), marker.getPopup(), options);
        });

        this.prevMarkerLists.set(
            type,
            list.map((e) => e.toObject())
        );
    }
}

export const centerX = 117.3;
export const centerY = 172.8;
export const scaleX = 0.02072;
export const scaleY = 0.0205;

export const customCRS = L.extend({}, L.CRS.Simple, {
    projection: L.Projection.LonLat,
    scale: function (zoom: number) {
        return Math.pow(2, zoom);
    },
    zoom: function (sc: number) {
        return Math.log(sc) / 0.6931471805599453;
    },
    distance: function (pos1: L.LatLng, pos2: L.LatLng) {
        var x_difference = pos2.lng - pos1.lng;
        var y_difference = pos2.lat - pos1.lat;
        return Math.sqrt(x_difference * x_difference + y_difference * y_difference);
    },
    transformation: new L.Transformation(scaleX, centerX, -scaleY, centerY),
    infinite: true,
});

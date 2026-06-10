import { CRS, extend, latLngBounds, Projection, Transformation, type LatLng } from 'leaflet';
import { tileLayers } from '~/types/livemap';

export const mapBackgroundColors = {
    postal: '#74aace',
    satellite: '#133e6b',
} as const;

const centerX = 117.3;
const centerY = 172.8;
const scaleX = 0.02072;
const scaleY = 0.0205;

export const mapBounds = latLngBounds([-4_000, -4_000], [8_000, 8_000]);
export const mapMaxBounds = latLngBounds([-9_000, -9_000], [11_000, 11_000]);

export const mapTileLayers = tileLayers;

export function getMapBackgroundColor(layer: string): string {
    return mapBackgroundColors[layer as keyof typeof mapBackgroundColors] ?? mapBackgroundColors.postal;
}

export const customMapCRS = extend({}, CRS.Simple, {
    projection: Projection.LonLat,
    scale: function (zoom: number): number {
        return Math.pow(2, zoom);
    },
    zoom: function (sc: number): number {
        return Math.log(sc) / 0.6931471805599453;
    },
    distance: function (pos1: LatLng, pos2: LatLng): number {
        const xDiff = pos2.lng - pos1.lng;
        const yDiff = pos2.lat - pos1.lat;
        return Math.sqrt(xDiff * xDiff + yDiff * yDiff);
    },
    transformation: new Transformation(scaleX, centerX, -scaleY, centerY),
    infinite: true,
});

export function getZoomOffset(zoom: number): number {
    switch (zoom) {
        case 1:
            return 1950;
        case 2:
            return 1450;
        case 3:
            return 1150;
        case 4:
            return 650;
        case 5:
            return 375;
        case 6:
            return 200;
        case 7:
            return 75;
        default:
            return 350;
    }
}

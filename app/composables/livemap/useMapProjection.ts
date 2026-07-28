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
    // GTA/FiveM coordinates are plain Cartesian coordinates, not geographic
    // latitude/longitude. We still use Leaflet's LonLat projection because it
    // preserves the values as x=lng and y=lat before the transformation below.
    projection: Projection.LonLat,
    scale: function (zoom: number): number {
        // Keep normal Leaflet zoom behavior: each zoom level doubles map pixels.
        return Math.pow(2, zoom);
    },
    zoom: function (sc: number): number {
        // Inverse of scale(). Leaflet passes an arbitrary scale value here.
        return Math.log(sc) / 0.6931471805599453;
    },
    distance: function (pos1: LatLng, pos2: LatLng): number {
        // Distance is measured in in-game units, so calculate Euclidean distance
        // directly from the unprojected CRS coordinates.
        const xDiff = pos2.lng - pos1.lng;
        const yDiff = pos2.lat - pos1.lat;
        return Math.sqrt(xDiff * xDiff + yDiff * yDiff);
    },
    // Affine transform from in-game coordinates to map pixels at zoom 0:
    //   mapX = x * scaleX + centerX
    //   mapY = y * -scaleY + centerY
    //
    // The negative Y scale flips GTA/FiveM's Y axis to match image tile space.
    // Leaflet multiplies both values by 2^zoom after this transform.
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

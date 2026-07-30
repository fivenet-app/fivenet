import type { ContextMenuItem } from '@nuxt/ui';
import type { LatLng, LatLngBoundsExpression } from 'leaflet';
import type { Perms } from '~~/gen/ts/perms';

export const tileLayers = [
    {
        key: 'postal',
        label: 'components.livemap.tile_layers.postal',
        url: '/images/livemap/tiles/postal/{z}/{x}/{y}.webp',
        options: {
            attribution: '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
            maxZoom: 7,
        },
    },
    {
        key: 'satellite',
        label: 'components.livemap.tile_layers.satellite',
        url: '/images/livemap/tiles/satellite/{z}/{x}/{y}.webp',
        options: {
            attribution: '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
            maxZoom: 7,
        },
    },
] as const;

type TileLayerItem = (typeof tileLayers)[number];

export type TileLayerKeys = TileLayerItem['key'];

export type Postal = {
    x: number;
    y: number;
    code: string;
};

export type LivemapContextMenuItem = Omit<ContextMenuItem, 'onSelect'> & {
    onSelect?: (latlng: LatLng) => Promise<void> | void;
    permission?: Perms;
};

export const overlayCayoPericoBounds: LatLngBoundsExpression = [
    [-6144.05, 3666.78],
    [-4144.05, 5723.78],
] as const;

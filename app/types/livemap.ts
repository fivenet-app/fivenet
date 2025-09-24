export const backgroundColorList = {
    postal: '#74aace',
    satelite: '#133e6b',
} as const;

export const tileLayers = [
    {
        key: 'postal',
        label: 'components.livemap.tile_layers.postal',
        url: '/images/livemap/tiles/postal/{z}/{x}/{y}.png',
        options: {
            attribution: '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
            maxZoom: 7,
        },
    },
    {
        key: 'satelite',
        label: 'components.livemap.tile_layers.satelite',
        url: '/images/livemap/tiles/satelite/{z}/{x}/{y}.png',
        options: {
            attribution: '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>',
            maxZoom: 7,
        },
    },
];

export type Postal = {
    x: number;
    y: number;
    code: string;
};

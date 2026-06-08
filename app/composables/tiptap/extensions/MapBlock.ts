/* eslint-disable @typescript-eslint/no-explicit-any */
import { mergeAttributes, Node } from '@tiptap/core';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import MapBlockNodeView from '~/components/partials/editor/MapBlockNodeView.vue';
import { tileLayers } from '~/types/livemap';

export interface MapBlockAttrs {
    x: number;
    y: number;
    zoom: number;
    postal?: string;
    layer?: string;
}

export interface MapBlockOptions {
    HTMLAttributes: Record<string, any>;
}

export const defaultMapBlockLayerKey = tileLayers[0]!.key;

type MapBlockAttrsInput = Partial<MapBlockAttrs> | Record<string, unknown> | undefined;

function readMapBlockAttr(attrs: MapBlockAttrsInput, key: keyof MapBlockAttrs): unknown {
    return (attrs as Record<string, unknown> | undefined)?.[key];
}

function parseMapBlockNumber(value: unknown, fallback: number): number {
    if (value === null || value === undefined || value === '') return fallback;

    const parsed = Number.parseFloat(String(value));
    return Number.isFinite(parsed) ? parsed : fallback;
}

export function normalizeMapBlockLayer(value: string | null | undefined): string {
    if (!value) return defaultMapBlockLayerKey;

    return tileLayers.some((layer) => layer.key === value) ? value : defaultMapBlockLayerKey;
}

export function normalizeMapBlockAttrs(attrs: MapBlockAttrsInput = {}): MapBlockAttrs {
    return {
        x: parseMapBlockNumber(readMapBlockAttr(attrs, 'x'), 0),
        y: parseMapBlockNumber(readMapBlockAttr(attrs, 'y'), 0),
        zoom: parseMapBlockNumber(readMapBlockAttr(attrs, 'zoom'), 2),
        postal: String(readMapBlockAttr(attrs, 'postal') ?? ''),
        layer: normalizeMapBlockLayer(String(readMapBlockAttr(attrs, 'layer') ?? '') || undefined),
    };
}

export function createMapBlockAttrs(attrs: MapBlockAttrsInput = {}): MapBlockAttrs {
    return normalizeMapBlockAttrs(attrs);
}

declare module '@tiptap/core' {
    interface Commands<ReturnType> {
        mapBlock: {
            insertMapBlock: (payload: MapBlockAttrs) => ReturnType;
        };
    }
}

function parseNumber(value: string | null | undefined, fallback: number): number {
    if (value === null || value === undefined || value === '') return fallback;

    const parsed = Number.parseFloat(value);
    return Number.isFinite(parsed) ? parsed : fallback;
}

function parseLayer(value: string | null | undefined): string {
    if (!value) return defaultMapBlockLayerKey;

    return tileLayers.some((layer) => layer.key === value) ? value : defaultMapBlockLayerKey;
}

export const MapBlock = Node.create<MapBlockOptions>({
    name: 'mapBlock',

    inline: true,
    group: 'inline',
    atom: true,
    draggable: true,
    selectable: true,
    isolating: true,

    addOptions() {
        return {
            HTMLAttributes: {},
        };
    },

    addAttributes() {
        return {
            x: {
                default: 0,
                parseHTML: (element) => parseNumber(element.getAttribute('data-map-x'), 0),
                renderHTML: (attributes) => ({ 'data-map-x': attributes.x }),
            },
            y: {
                default: 0,
                parseHTML: (element) => parseNumber(element.getAttribute('data-map-y'), 0),
                renderHTML: (attributes) => ({ 'data-map-y': attributes.y }),
            },
            zoom: {
                default: 2,
                parseHTML: (element) => parseNumber(element.getAttribute('data-map-zoom'), 2),
                renderHTML: (attributes) => ({ 'data-map-zoom': attributes.zoom }),
            },
            layer: {
                default: defaultMapBlockLayerKey,
                parseHTML: (element) => parseLayer(element.getAttribute('data-map-layer')),
                renderHTML: (attributes) => ({ 'data-map-layer': attributes.layer || defaultMapBlockLayerKey }),
            },
            postal: {
                default: '',
                parseHTML: (element) => element.getAttribute('data-map-postal') ?? '',
                renderHTML: (attributes) => {
                    if (!attributes.postal) return {};
                    return { 'data-map-postal': attributes.postal };
                },
            },
        };
    },

    parseHTML() {
        return [{ tag: 'figure[data-embed="map"]' }, { tag: 'div[data-embed="map"]' }, { tag: 'span[data-embed="map"]' }];
    },

    renderHTML({ HTMLAttributes }) {
        const attrs = {
            ...HTMLAttributes,
            class: 'map-block inline-flex align-middle',
            'data-embed': 'map',
        };

        return ['span', mergeAttributes(this.options.HTMLAttributes, attrs)];
    },

    addCommands() {
        return {
            insertMapBlock:
                (payload) =>
                ({ commands }) =>
                    commands.insertContent({
                        type: this.name,
                        attrs: {
                            x: payload.x,
                            y: payload.y,
                            zoom: payload.zoom,
                            postal: payload.postal ?? '',
                            layer: payload.layer ?? defaultMapBlockLayerKey,
                        },
                    }),
        };
    },

    addNodeView() {
        return VueNodeViewRenderer(MapBlockNodeView);
    },
});

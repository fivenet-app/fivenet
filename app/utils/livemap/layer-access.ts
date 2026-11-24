import type { LivemapLayer } from '~/stores/settings';

export type LivemapLayerAccess = {
    can: (perm: NonNullable<LivemapLayer['perm']>) => boolean;
    attr: (perm: NonNullable<LivemapLayer['perm']>, key: string, value: string) => boolean;
};

/**
 * Determines whether a layer may be exposed in either the layer controls or
 * the map itself. Attribute access supplements, rather than replaces, the
 * layer's base permission.
 */
export function canAccessLivemapLayer(layer: LivemapLayer, access: LivemapLayerAccess): boolean {
    if (layer.perm === undefined) return true;
    if (!access.can(layer.perm)) return false;
    if (layer.attr === undefined) return true;

    return access.attr(layer.perm, layer.attr.key, layer.attr.val);
}

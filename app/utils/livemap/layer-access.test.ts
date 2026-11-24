import { describe, expect, it } from 'vitest';
import type { LivemapLayer } from '~/stores/settings';
import { canAccessLivemapLayer } from './layer-access';

const perm = 'livemap.LivemapService/Stream' as const;

const layer = (overrides: Partial<LivemapLayer> = {}): LivemapLayer => ({
    key: 'test',
    label: 'Test',
    category: 'test',
    ...overrides,
});

describe('canAccessLivemapLayer', () => {
    it('allows layers without permission requirements', () => {
        expect(canAccessLivemapLayer(layer(), { can: () => false, attr: () => false })).toBe(true);
    });

    it('allows a layer with its base permission', () => {
        expect(canAccessLivemapLayer(layer({ perm }), { can: () => true, attr: () => false })).toBe(true);
    });

    it('denies an attribute-scoped layer when its base permission is missing', () => {
        expect(
            canAccessLivemapLayer(layer({ perm, attr: { key: 'Markers', val: 'police' } }), {
                can: () => false,
                attr: () => true,
            }),
        ).toBe(false);
    });

    it('denies an attribute-scoped layer when its attribute is missing', () => {
        expect(
            canAccessLivemapLayer(layer({ perm, attr: { key: 'Players', val: 'police' } }), {
                can: () => true,
                attr: () => false,
            }),
        ).toBe(false);
    });

    it('allows an attribute-scoped layer only when both checks pass', () => {
        expect(
            canAccessLivemapLayer(layer({ perm, attr: { key: 'Markers', val: 'police' } }), {
                can: () => true,
                attr: () => true,
            }),
        ).toBe(true);
    });
});

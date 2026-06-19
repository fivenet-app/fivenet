import { describe, expect, it } from 'vitest';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';
import { resolveWikiPageMovePayload } from './reorder';

function page(id: number): PageShort {
    return {
        id,
    } as PageShort;
}

describe('resolveWikiPageMovePayload', () => {
    it('uses beforeId when moving a page upward', () => {
        const siblings = [page(348), page(319)];

        expect(resolveWikiPageMovePayload(siblings, 1, 0)).toEqual({
            pageId: 348,
            beforeId: 319,
        });
    });

    it('uses afterId when moving a page downward', () => {
        const siblings = [page(319), page(348)];

        expect(resolveWikiPageMovePayload(siblings, 0, 1)).toEqual({
            pageId: 348,
            afterId: 319,
        });
    });

    it('returns undefined when the page does not move', () => {
        const siblings = [page(319), page(348)];

        expect(resolveWikiPageMovePayload(siblings, 1, 1)).toBeUndefined();
    });
});

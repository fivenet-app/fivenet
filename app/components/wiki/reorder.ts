import { resolveNeighborMovePayload } from '~/utils/reorder';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';
import { sameWikiMoveGroup } from './helpers';

export type WikiPageMovePayload = {
    pageId: number;
    beforeId?: number;
    afterId?: number;
};

type WikiPageDragEvent = {
    dragged?: HTMLElement;
    related?: HTMLElement;
    oldIndex?: number;
    newIndex?: number;
};

function getPageIdFromElement(element: HTMLElement | undefined): number | undefined {
    if (!element) return undefined;

    const pageId = element.dataset.pageId;
    if (!pageId) return undefined;

    const parsed = Number(pageId);
    return Number.isFinite(parsed) ? parsed : undefined;
}

function getPageFromElement(element: HTMLElement | undefined, siblings: PageShort[]): PageShort | undefined {
    const pageId = getPageIdFromElement(element);
    if (pageId === undefined) return undefined;

    return siblings.find((page) => page.id === pageId);
}

export function canReorderWikiPages(siblings: PageShort[], event: WikiPageDragEvent): boolean {
    const dragged = getPageFromElement(event.dragged, siblings);
    const related = getPageFromElement(event.related, siblings);

    if (!dragged || !related) return true;

    return sameWikiMoveGroup(dragged, related);
}

export function resolveWikiPageMovePayload(
    siblings: PageShort[],
    oldIndex: number | undefined,
    newIndex: number | undefined,
): WikiPageMovePayload | undefined {
    const payload = resolveNeighborMovePayload(siblings, oldIndex, newIndex);
    if (!payload) return undefined;

    return {
        pageId: payload.id,
        beforeId: payload.beforeId,
        afterId: payload.afterId,
    };
}

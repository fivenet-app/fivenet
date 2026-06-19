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
    if (oldIndex === undefined || newIndex === undefined || oldIndex === newIndex) return undefined;

    const page = siblings[newIndex];
    if (!page) return undefined;

    if (newIndex < oldIndex) {
        const beforeId = siblings[newIndex + 1]?.id;
        return beforeId ? { pageId: page.id, beforeId } : undefined;
    }

    const afterId = siblings[newIndex - 1]?.id;
    if (!afterId) return undefined;

    return { pageId: page.id, afterId };
}

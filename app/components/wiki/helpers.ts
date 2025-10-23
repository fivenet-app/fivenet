// Pageument Activity

import type { AccessLevel, PageAccess } from '~~/gen/ts/resources/wiki/access';
import { PageActivityType } from '~~/gen/ts/resources/wiki/activity';
import { type Page, PageShort } from '~~/gen/ts/resources/wiki/page';

export function getPageAtivityIcon(activityType: PageActivityType): string {
    switch (activityType) {
        // Base
        case PageActivityType.CREATED:
            return 'i-mdi-new-box';
        case PageActivityType.UPDATED:
            return 'i-mdi-update';
        case PageActivityType.ACCESS_UPDATED:
            return 'i-mdi-lock-check';
        case PageActivityType.OWNER_CHANGED:
            return 'i-mdi-file-account';
        case PageActivityType.DELETED:
            return 'i-mdi-delete-circle';
        case PageActivityType.DRAFT_TOGGLED:
            return 'i-mdi-file-document-edit-outline';

        default:
            return 'i-mdi-help';
    }
}

export function checkPageAccess(
    access: PageAccess | undefined,
    creator: UserLike | undefined,
    level: AccessLevel,
    creatorJob?: string,
): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) return true;

    if (activeChar.value === null) return false;

    return checkAccess(activeChar.value, access, creator, level, creatorJob);
}

export function pageToURL(page: PageShort | Page, fullUrl: boolean = false): string {
    const base = fullUrl ? `${window.location.protocol}//${window.location.hostname}` : '';

    if (PageShort.is(page)) {
        return `${base}/wiki/${page.job}/${page.id}/${page?.slug ?? ''}`;
    } else {
        return `${base}/wiki/${page.job}/${page.id}/${page?.meta?.slug ?? ''}`;
    }
}

// Pageument Activity

import type { AccessLevel, PageAccess } from '~~/gen/ts/resources/wiki/access';
import { PageActivityType } from '~~/gen/ts/resources/wiki/activity';

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

        default:
            return 'i-mdi-help';
    }
}

export function checkPageAccess(access: PageAccess | undefined, creator: UserLike | undefined, level: AccessLevel) {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    return checkAccess(activeChar.value, access, creator, level);
}

// Pageument Activity

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

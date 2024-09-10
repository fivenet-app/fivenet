import { useAuthStore } from '~/store/auth';
import type { AccessLevel, CalendarAccess } from '~~/gen/ts/resources/calendar/access';
import type { UserShort } from '~~/gen/ts/resources/users/users';

export function checkCalendarAccess(
    access: CalendarAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    const activeChar = authStore.activeChar;
    if (activeChar === null) {
        return false;
    }

    return checkAccess(activeChar, access, creator, level);
}

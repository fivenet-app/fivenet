import type { AccessLevel, CalendarAccess } from '~~/gen/ts/resources/calendar/access/access';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';

export function checkCalendarAccess(
    access: CalendarAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    return checkAccess(activeChar.value, access, creator, level);
}

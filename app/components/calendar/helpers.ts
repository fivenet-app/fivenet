import { useAuthStore } from '~/store/auth';
import type { AccessLevel, CalendarAccess } from '~~/gen/ts/resources/calendar/access';
import type { UserShort } from '~~/gen/ts/resources/users/users';

export function checkCalendarAccess(
    calendarAccess: CalendarAccess | undefined,
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

    if (calendarAccess === undefined) {
        return false;
    }

    if (creator !== undefined && activeChar.userId === creator.userId) {
        return true;
    }

    const ju = calendarAccess.users.find((ua) => ua.userId === activeChar.userId && level <= ua.access);
    if (ju !== undefined) {
        return true;
    }

    const ja = calendarAccess.jobs.find(
        (ja) => ja.job === activeChar.job && ja.minimumGrade <= activeChar.jobGrade && level <= ja.access,
    );
    if (ja !== undefined) {
        return true;
    }

    return false;
}

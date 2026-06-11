import type { ButtonProps } from '@nuxt/ui';
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import type { UserLike } from '~/utils/strings';
import type { AccessLevel, CalendarAccess } from '~~/gen/ts/resources/calendar/access/access';
import { CalendarSystemKind, type Calendar } from '~~/gen/ts/resources/calendar/calendar';
import { CalendarEntryOccurrenceKind, type CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';

export type CalendarEntryAttributeData = CalendarEntry & {
    color: ButtonProps['color'];
    icon?: string;
    isPast: boolean;
    multiDay: boolean;
    ongoing: boolean;
    time: string;
    timeEnd?: string;
};

export type CalendarEntryAttribute = {
    key: string;
    customData: CalendarEntryAttributeData;
    dates: DateRangeSource;
};

export function isCalendarCreator(activeChar: UserLike, creator?: UserShort, calendarJob?: string): boolean {
    return creator !== undefined && calendarJob === undefined && activeChar.userId === creator.userId;
}

export function checkCalendarAccess(
    access: CalendarAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    calendarJob?: string,
    creatorJob?: string,
): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) return true;

    if (activeChar.value === null) return false;

    if (isCalendarCreator(activeChar.value, creator, calendarJob)) return true;

    return checkAccess(activeChar.value, access, creator, level, creatorJob);
}

export function isSystemManagedCalendar(calendar?: Pick<Calendar, 'systemKind'> | undefined): boolean {
    return calendar?.systemKind !== undefined && calendar?.systemKind > CalendarSystemKind.UNSPECIFIED;
}

export function isBirthdayEntry(entry?: Pick<CalendarEntry, 'occurrence'> | undefined): boolean {
    return entry?.occurrence?.kind === CalendarEntryOccurrenceKind.BIRTHDAY;
}

export function isSystemManagedCalendarEntry(
    calendar?: Pick<Calendar, 'systemKind'> | undefined,
    entry?: Pick<CalendarEntry, 'occurrence'> | undefined,
): boolean {
    return isSystemManagedCalendar(calendar) || isBirthdayEntry(entry);
}

import type { ButtonProps } from '@nuxt/ui';
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import type { AccessLevel, CalendarAccess } from '~~/gen/ts/resources/calendar/access/access';
import type { Calendar } from '~~/gen/ts/resources/calendar/calendar';
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

export function checkCalendarAccess(
    access: CalendarAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    creatorJob?: string,
): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) return true;

    if (activeChar.value === null) return false;

    return checkAccess(activeChar.value, access, creator, level, creatorJob);
}

export function isSystemManagedCalendar(calendar?: Pick<Calendar, 'systemKind'> | undefined): boolean {
    return calendar?.systemKind !== undefined;
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

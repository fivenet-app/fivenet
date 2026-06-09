import { CalendarEntryOccurrenceKind, type CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import { toDate } from './time';

type BirthdayOccurrenceDate = {
    year: number;
    month: number;
    day: number;
};

export function parseBirthdayOccurrenceKey(key: string): BirthdayOccurrenceDate | undefined {
    const [kind, _calendarId, _userId, year, month, day] = key.split(':');

    if (kind !== 'birthday' || !_calendarId || !_userId || !year || !month || !day) return undefined;

    const parsedYear = Number(year);
    const parsedMonth = Number(month);
    const parsedDay = Number(day);

    if (!Number.isFinite(parsedYear) || !Number.isFinite(parsedMonth) || !Number.isFinite(parsedDay)) return undefined;

    return {
        year: parsedYear,
        month: parsedMonth,
        day: parsedDay,
    };
}

export function getCalendarEntryDisplayStartDate(entry: CalendarEntry): Date {
    const occurrence = entry.occurrence;

    if (occurrence?.kind === CalendarEntryOccurrenceKind.BIRTHDAY) {
        const parsed = occurrence.key ? parseBirthdayOccurrenceKey(occurrence.key) : undefined;
        if (parsed) {
            return new Date(parsed.year, parsed.month - 1, parsed.day, 12, 0, 0, 0);
        }
    }

    return toDate(entry.startTime);
}

export function getCalendarEntryDisplayEndDate(entry: CalendarEntry): Date | undefined {
    if (entry.occurrence?.kind === CalendarEntryOccurrenceKind.BIRTHDAY) {
        return undefined;
    }

    return entry.endTime ? toDate(entry.endTime) : undefined;
}

export function getCalendarEntryDisplayRangeEndDate(entry: CalendarEntry): Date | undefined {
    const end = getCalendarEntryDisplayEndDate(entry);
    if (end) return end;

    if (entry.occurrence?.allDay) {
        return getCalendarEntryDisplayStartDate(entry);
    }

    return undefined;
}

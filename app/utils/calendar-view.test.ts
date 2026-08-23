import { describe, expect, it } from 'vitest';
import { toTimestamp } from '~/utils/time';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import {
    formatCalendarRouteDate,
    getCalendarEntryIcon,
    getCalendarEntryRangeEnd,
    layoutAllDayEntries,
    layoutDayEntries,
    parseCalendarRouteDate,
    routeForCalendarView,
    viewDateRange,
} from './calendar-view';

function makeEntry(id: number, start: string, end?: string, allDay = false): CalendarEntry {
    return {
        id,
        calendarId: 1,
        title: `Entry ${id}`,
        startTime: toTimestamp(new Date(start))!,
        endTime: end ? toTimestamp(new Date(end)) : undefined,
        occurrence: allDay
            ? {
                  allDay: true,
              }
            : undefined,
    } as CalendarEntry;
}

describe('calendar-view helpers', () => {
    it('parses and formats route dates', () => {
        const date = parseCalendarRouteDate('2026-08-17');

        expect(date).toBeTruthy();
        expect(formatCalendarRouteDate(date!)).toBe('2026-08-17');
        expect(routeForCalendarView('month', date!)).toBe('/calendar/month/2026-08-17');
    });

    it('returns a week range anchored on monday', () => {
        const range = viewDateRange('week', new Date('2026-08-17T12:00:00'));

        expect(formatCalendarRouteDate(range.start)).toBe('2026-08-17');
        expect(formatCalendarRouteDate(range.end)).toBe('2026-08-24');
    });

    it('packs overlapping day events into columns', () => {
        const positioned = layoutDayEntries(
            [
                makeEntry(1, '2026-08-17T09:00:00', '2026-08-17T10:00:00'),
                makeEntry(2, '2026-08-17T09:30:00', '2026-08-17T10:30:00'),
            ],
            new Date('2026-08-17T00:00:00'),
        );

        expect(positioned).toHaveLength(2);
        expect(positioned[0]?.width).toBeLessThan(100);
        expect(positioned[1]?.left).toBeGreaterThanOrEqual(0);
    });

    it('keeps nearby non-overlapping day events full width', () => {
        const positioned = layoutDayEntries(
            [
                makeEntry(1, '2026-08-17T09:00:00', '2026-08-17T09:05:00'),
                makeEntry(2, '2026-08-17T09:20:00', '2026-08-17T09:25:00'),
            ],
            new Date('2026-08-17T00:00:00'),
        );

        expect(positioned).toHaveLength(2);
        expect(positioned[0]?.width).toBe(100);
        expect(positioned[1]?.width).toBe(100);
    });

    it('clips multi-day timed events to the visible day segment', () => {
        const entry = makeEntry(1, '2026-08-17T01:00:00', '2026-08-18T14:00:00');

        const startDay = layoutDayEntries([entry], new Date('2026-08-17T00:00:00'));
        const middleDay = layoutDayEntries([entry], new Date('2026-08-18T00:00:00'));

        expect(startDay).toHaveLength(1);
        expect(middleDay).toHaveLength(1);
        expect(startDay[0]?.top).toBeCloseTo(60 * (64 / 60));
        expect(startDay[0]?.height).toBeCloseTo(23 * 64);
        expect(middleDay[0]?.top).toBe(0);
        expect(middleDay[0]?.height).toBeCloseTo(14 * 64);
    });

    it('widens solitary events in week mode even when they belong to an overlap chain', () => {
        const positioned = layoutDayEntries(
            [
                makeEntry(1, '2026-08-17T09:00:00', '2026-08-17T11:30:00'),
                makeEntry(2, '2026-08-17T10:00:00', '2026-08-17T10:30:00'),
                makeEntry(3, '2026-08-17T11:45:00', '2026-08-17T12:15:00'),
            ],
            new Date('2026-08-17T00:00:00'),
            { preferWideSolitaryEvents: true },
        );

        expect(positioned).toHaveLength(3);
        expect(positioned[2]?.width).toBe(100);
    });

    it('packs all-day events into lanes', () => {
        const positioned = layoutAllDayEntries(
            [
                makeEntry(1, '2026-08-17T00:00:00', '2026-08-18T00:00:00', true),
                makeEntry(2, '2026-08-17T00:00:00', '2026-08-18T00:00:00', true),
            ],
            [new Date('2026-08-17T00:00:00')],
        );

        expect(positioned).toHaveLength(2);
        expect(positioned[0]?.lane).toBe(0);
        expect(positioned[1]?.lane).toBe(1);
    });

    it('spans multi-day all-day events across every day in the range', () => {
        const positioned = layoutAllDayEntries(
            [makeEntry(1, '2026-08-17T00:00:00', '2026-08-19T00:00:00', true)],
            [
                new Date('2026-08-17T00:00:00'),
                new Date('2026-08-18T00:00:00'),
                new Date('2026-08-19T00:00:00'),
                new Date('2026-08-20T00:00:00'),
            ],
        );

        expect(positioned).toHaveLength(1);
        expect(positioned[0]?.colStart).toBe(0);
        expect(positioned[0]?.colSpan).toBe(3);
    });

    it('treats all-day entries without an explicit end as spanning one day', () => {
        const entry = makeEntry(1, '2026-08-17T00:00:00', undefined, true);
        const end = getCalendarEntryRangeEnd(entry);

        expect(end.getFullYear()).toBe(2026);
        expect(end.getMonth()).toBe(7);
        expect(end.getDate()).toBe(18);
    });

    it('prefers a custom entry icon over derived fallback icons', () => {
        const entry = makeEntry(1, '2026-08-17T10:00:00', '2026-08-17T11:00:00');
        entry.icon = 'CalendarStarIcon';

        expect(getCalendarEntryIcon(entry)).toBe('i-mdi-calendar-star');
    });
});

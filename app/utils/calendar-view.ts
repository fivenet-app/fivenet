import type { ButtonProps } from '@nuxt/ui';
import {
    addDays,
    addHours,
    addMinutes,
    differenceInCalendarDays,
    eachMonthOfInterval,
    endOfMonth,
    format,
    isAfter,
    isBefore,
    isSameDay,
    startOfDay,
    startOfMonth,
    startOfWeek,
} from 'date-fns';
import { isValidCalendarEntryRecurring } from '~/components/calendar/helpers';
import {
    getCalendarEntryDisplayEndDate,
    getCalendarEntryDisplayRangeEndDate,
    getCalendarEntryDisplayStartDate,
    isCalendarEntryAllDay,
} from '~/utils/calendar';
import { dateToDateString } from '~/utils/time';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';

export type CalendarView = 'day' | 'week' | 'month' | 'summary';

export type CalendarEntryPosition = {
    entry: CalendarEntry;
    top: number;
    height: number;
    left: number;
    width: number;
};

export type CalendarAllDayEntryPosition = {
    entry: CalendarEntry;
    colStart: number;
    colSpan: number;
    lane: number;
};

export type LayoutDayEntriesOptions = {
    preferWideSolitaryEvents?: boolean;
};

export const HOUR_HEIGHT = 64;
export const PX_PER_MINUTE = HOUR_HEIGHT / 60;
export const MIN_EVENT_MINUTES = 30;

export function getCalendarCreateRangeFromClick(day: Date, event: MouseEvent): { startTime: Date; endTime: Date } | undefined {
    const target = event.currentTarget as HTMLElement | null;
    if (!target) return undefined;

    const rect = target.getBoundingClientRect();
    const relativeY = Math.max(0, Math.min(event.clientY - rect.top, 24 * HOUR_HEIGHT - 1));
    const hour = Math.min(23, Math.floor(relativeY / HOUR_HEIGHT));
    const startTime = new Date(day);
    startTime.setHours(hour, 0, 0, 0);

    return {
        startTime,
        endTime: addHours(startTime, 1),
    };
}

export function isCalendarView(view: string): view is CalendarView {
    return view === 'day' || view === 'week' || view === 'month' || view === 'summary';
}

export function normalizeCalendarView(view: string | undefined): CalendarView {
    return view && isCalendarView(view) ? view : 'month';
}

export function parseCalendarRouteDate(value: string | undefined): Date | undefined {
    if (!value) return undefined;

    const parsed = new Date(`${value}T00:00:00`);
    return Number.isNaN(parsed.getTime()) ? undefined : parsed;
}

export function formatCalendarRouteDate(date: Date): string {
    return dateToDateString(date);
}

export function routeForCalendarView(view: CalendarView, date: Date): string {
    return `/calendar/${view}/${formatCalendarRouteDate(date)}`;
}

export function viewDateRange(view: CalendarView, date: Date): { start: Date; end: Date } {
    if (view === 'day') {
        const start = startOfDay(date);
        return { start, end: addDays(start, 1) };
    }

    if (view === 'week') {
        const start = startOfWeek(date, { weekStartsOn: 1 });
        return { start, end: addDays(start, 7) };
    }

    const start = startOfWeek(startOfMonth(date), { weekStartsOn: 1 });
    return { start, end: addDays(start, 42) };
}

export function fetchMonthsForRange(range: { start: Date; end: Date }): Date[] {
    return eachMonthOfInterval({
        start: startOfMonth(range.start),
        end: endOfMonth(addDays(range.end, -1)),
    });
}

export function getVisibleCalendarEntryIds(entries: CalendarEntry[]): string[] {
    return entries.map((entry) => entry.occurrence?.key ?? String(entry.id));
}

export function getCalendarEntryColor(entry: CalendarEntry): ButtonProps['color'] {
    return (entry.calendar?.color as ButtonProps['color'] | undefined) ?? 'primary';
}

export function getCalendarEntryIcon(entry: CalendarEntry): string | undefined {
    if (entry.calendar?.systemKind) return 'i-mdi-badge-account-horizontal-outline';
    if (entry.deletedAt) return 'i-mdi-delete';
    if (isValidCalendarEntryRecurring(entry.recurring)) return 'i-mdi-repeat';
    if (isCalendarEntryAllDay(entry)) return 'i-mdi-calendar';
    return undefined;
}

export function getCalendarEntryTimeLabel(entry: CalendarEntry): string {
    const start = getCalendarEntryDisplayStartDate(entry);
    const end = getCalendarEntryDisplayEndDate(entry);

    if (isCalendarEntryAllDay(entry)) {
        if (end) {
            if (!isSameDay(start, end)) {
                return `${format(start, 'P')} - ${format(end, 'P')}`;
            }
        }

        return format(start, 'P');
    }

    if (!end) {
        return format(start, 'p');
    }

    if (isSameDay(start, end)) {
        return `${format(start, 'p')} - ${format(end, 'p')}`;
    }

    return `${format(start, 'P p')} - ${format(end, 'P p')}`;
}

export function getCalendarEntryRangeEnd(entry: CalendarEntry): Date {
    const rangeEnd = getCalendarEntryDisplayRangeEndDate(entry);
    if (rangeEnd) {
        return isCalendarEntryAllDay(entry) ? getCalendarEntryComparisonEndDate(entry) : rangeEnd;
    }

    const start = getCalendarEntryDisplayStartDate(entry);
    return addMinutes(start, MIN_EVENT_MINUTES);
}

export function isCalendarEntryPast(entry: CalendarEntry, now = new Date()): boolean {
    const start = getCalendarEntryDisplayStartDate(entry);

    if (isCalendarEntryAllDay(entry)) {
        return isAfter(now, startOfDay(start));
    }

    const end = getCalendarEntryDisplayEndDate(entry);
    return isAfter(now, end ?? start);
}

export function isCalendarEntryOngoing(entry: CalendarEntry, now = new Date()): boolean {
    const start = getCalendarEntryDisplayStartDate(entry);

    if (isCalendarEntryAllDay(entry)) {
        return isBefore(start, now) && isBefore(now, getCalendarEntryRangeEnd(entry));
    }

    const end = getCalendarEntryDisplayEndDate(entry);
    return isBefore(start, now) && !!end && isAfter(end, now);
}

export function layoutDayEntries(
    entries: CalendarEntry[],
    day: Date,
    options: LayoutDayEntriesOptions = {},
): CalendarEntryPosition[] {
    const dayStart = startOfDay(day);
    const dayEnd = addDays(dayStart, 1);

    const items = entries
        .map((entry) => {
            const start = getCalendarEntryDisplayStartDate(entry);
            const end = getCalendarEntryRangeEnd(entry);
            const segmentStart = start < dayStart ? dayStart : start;
            const segmentEnd = end > dayEnd ? dayEnd : end;
            const startMin = Math.max(0, Math.floor((segmentStart.getTime() - dayStart.getTime()) / 60000));
            const actualEndMin = Math.min(24 * 60, Math.ceil((segmentEnd.getTime() - dayStart.getTime()) / 60000));

            return {
                entry,
                startMin,
                actualEndMin,
                renderEndMin: Math.max(actualEndMin, startMin + MIN_EVENT_MINUTES),
            };
        })
        .filter((item) => {
            const start = getCalendarEntryDisplayStartDate(item.entry);
            const end = getCalendarEntryRangeEnd(item.entry);
            return start < dayEnd && end > dayStart && item.actualEndMin > item.startMin;
        })
        .sort((a, b) => a.startMin - b.startMin || b.actualEndMin - a.actualEndMin);

    const positioned: CalendarEntryPosition[] = [];
    let cluster: typeof items = [];
    let clusterEnd = -Infinity;

    function flush(): void {
        const columns: number[] = [];
        const assigned = cluster.map((item) => {
            let index = columns.findIndex((end) => end <= item.startMin);
            if (index === -1) index = columns.length;
            columns[index] = item.actualEndMin;
            return index;
        });

        const overlapsRange = (leftStart: number, leftEnd: number, rightStart: number, rightEnd: number): boolean =>
            leftStart < rightEnd && leftEnd > rightStart;

        for (const [index, item] of cluster.entries()) {
            const ownColumn = assigned[index] ?? 0;
            let leftColumn = ownColumn;
            let rightColumn = ownColumn;
            const blockedColumns = new Set<number>();

            if (options.preferWideSolitaryEvents) {
                for (const [otherIndex, other] of cluster.entries()) {
                    if (otherIndex === index) continue;
                    if (!overlapsRange(item.startMin, item.actualEndMin, other.startMin, other.actualEndMin)) continue;

                    const otherColumn = assigned[otherIndex];
                    if (otherColumn !== undefined) {
                        blockedColumns.add(otherColumn);
                    }
                }

                while (leftColumn > 0 && !blockedColumns.has(leftColumn - 1)) {
                    leftColumn--;
                }

                while (rightColumn + 1 < columns.length && !blockedColumns.has(rightColumn + 1)) {
                    rightColumn++;
                }
            }

            const widthColumns = Math.max(1, rightColumn - leftColumn + 1);
            positioned.push({
                entry: item.entry,
                top: item.startMin * PX_PER_MINUTE,
                height: (item.renderEndMin - item.startMin) * PX_PER_MINUTE,
                left: (leftColumn / columns.length) * 100,
                width: (100 / columns.length) * widthColumns,
            });
        }
    }

    for (const item of items) {
        if (cluster.length && item.startMin >= clusterEnd) {
            flush();
            cluster = [];
            clusterEnd = -Infinity;
        }

        cluster.push(item);
        clusterEnd = Math.max(clusterEnd, item.actualEndMin);
    }

    if (cluster.length) {
        flush();
    }

    return positioned;
}

export function layoutAllDayEntries(entries: CalendarEntry[], days: Date[]): CalendarAllDayEntryPosition[] {
    const first = days[0];
    if (!first) return [];
    const firstDay = startOfDay(first);

    const items = entries
        .map((entry) => {
            const start = getCalendarEntryDisplayStartDate(entry);
            const end = getCalendarEntryRangeEnd(entry);
            const startDay = startOfDay(start);
            const endDayExclusive = startOfDay(end);
            const colStart = Math.max(0, differenceInCalendarDays(startDay, firstDay));
            const colEnd = Math.min(days.length, differenceInCalendarDays(endDayExclusive, firstDay));
            return { entry, colStart, colSpan: colEnd - colStart };
        })
        .filter((item) => item.colSpan > 0)
        .sort((a, b) => a.colStart - b.colStart || b.colSpan - a.colSpan);

    const lanes: number[] = [];

    return items.map((item) => {
        let lane = lanes.findIndex((end) => end <= item.colStart);
        if (lane === -1) lane = lanes.length;
        lanes[lane] = item.colStart + item.colSpan;
        return { ...item, lane };
    });
}

export function groupEntriesByDay(entries: CalendarEntry[]): Map<string, CalendarEntry[]> {
    const groups = new Map<string, CalendarEntry[]>();

    for (const entry of entries) {
        const start = startOfDay(getCalendarEntryDisplayStartDate(entry));
        const end = getCalendarEntryRangeEnd(entry);
        for (let day = start; day < end; day = addDays(day, 1)) {
            const key = dateToDateString(day);
            const list = groups.get(key) ?? [];
            list.push(entry);
            groups.set(key, list);
        }
    }

    return groups;
}

export function rangeTitle(view: CalendarView, date: Date): string {
    if (view === 'day') {
        return format(date, 'PPPP');
    }

    if (view === 'week') {
        const start = startOfWeek(date, { weekStartsOn: 1 });
        const end = addDays(start, 6);

        if (start.getMonth() === end.getMonth()) {
            return `${format(start, 'MMMM d')} - ${format(end, 'd, yyyy')}`;
        }

        return `${format(start, 'MMM d')} - ${format(end, 'MMM d, yyyy')}`;
    }

    return format(date, 'MMMM yyyy');
}

import { describe, expect, it } from 'vitest';
import { dateToDateString, getWeekNumber, stringToDate, toDate, toDatetimeLocal, toTimestamp } from './time';

describe('toDate', () => {
    it('should return the current date if input is undefined', () => {
        const result = toDate(undefined);
        expect(result).toBeInstanceOf(Date);
    });

    it('should return a corrected date if correction is provided', () => {
        const timestamp = { timestamp: { seconds: 1234567890, nanos: 0 } };
        const correction = 1000;
        const result = toDate(timestamp, correction);
        expect(result.getTime()).toBeLessThan(Date.now());
    });
});

describe('stringToDate', () => {
    it('should parse a valid date string into a Date object', () => {
        const input = '2025-11-25T12:00:00Z';
        const result = stringToDate(input);
        const output = '2025-11-25T12:00:00.000Z';
        expect(result.toISOString()).toBe(output);
    });
});

describe('toTimestamp', () => {
    it('should return undefined if input is undefined', () => {
        const result = toTimestamp();
        expect(result).toBeUndefined();
    });

    it('should convert a Date object to a Timestamp object', () => {
        const date = new Date('2025-11-25T12:00:00Z');
        const result = toTimestamp(date);
        expect(result).toEqual({ timestamp: { seconds: Math.floor(date.getTime() / 1000), nanos: 0 } });
    });
});

describe('toDatetimeLocal', () => {
    it('should convert a Date object to a datetime-local string', () => {
        const date = new Date(Date.UTC(2025, 10, 25, 12, 0, 0)); // November is month 10 (0-based)
        const result = toDatetimeLocal(date);

        // Dynamically calculate expected local time string
        const localYear = date.getFullYear();
        const localMonth = String(date.getMonth() + 1).padStart(2, '0');
        const localDay = String(date.getDate()).padStart(2, '0');
        const localHours = String(date.getHours()).padStart(2, '0');
        const localMinutes = String(date.getMinutes()).padStart(2, '0');
        const expected = `${localYear}-${localMonth}-${localDay}T${localHours}:${localMinutes}`;

        expect(result).toBe(expected);
    });
});

describe('dateToDateString', () => {
    it('should convert a Date object to a date string in YYYY-MM-DD format', () => {
        const date = new Date('2025-11-25');
        const result = dateToDateString(date);
        expect(result).toBe('2025-11-25');
    });
});

describe('getWeekNumber', () => {
    it('should return the correct ISO week number for a given date', () => {
        const date = new Date('2025-11-25');
        const result = getWeekNumber(date);
        expect(result).toBe(48);
    });

    it('should return the correct ISO week number for a given date (end of year)', () => {
        const date = new Date('2025-12-31');
        const result = getWeekNumber(date);
        expect(result).toBe(1);
    });
});

import { describe, expect, it } from 'vitest';
import { calculateOffset } from './pagination';

describe('calculateOffset', () => {
    it('should return 0 if pagination is undefined', () => {
        const page = 1;
        const result = calculateOffset(page);
        expect(result).toBe(0);
    });

    it('should calculate the correct offset when pagination is provided', () => {
        const page = 3;
        const pagination = { pageSize: 10, totalCount: 100, offset: 0, end: 10 };
        const result = calculateOffset(page, pagination);
        expect(result).toBe(20);
    });

    it('should return 0 for the first page', () => {
        const page = 1;
        const pagination = { pageSize: 10, totalCount: 100, offset: 0, end: 10 };
        const result = calculateOffset(page, pagination);
        expect(result).toBe(0);
    });
});

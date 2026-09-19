import { describe, expect, it } from 'vitest';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';
import { createExamNavigatorEntries, getExamNavigatorTargetIndex } from './examNavigator';

function question(id: number, kind: string): ExamQuestion {
    return { id, data: { data: { oneofKind: kind } } } as ExamQuestion;
}

describe('exam navigator helpers', () => {
    it('numbers only answerable questions', () => {
        const entries = createExamNavigatorEntries([
            question(1, 'separator'),
            question(2, 'yesno'),
            question(3, 'image'),
            question(4, 'freeText'),
        ]);

        expect(entries.map((entry) => [entry.kind, entry.answerable, entry.number])).toEqual([
            ['separator', false, undefined],
            ['yesno', true, 1],
            ['image', false, undefined],
            ['freeText', true, 2],
        ]);
    });

    it('calculates keyboard navigation targets', () => {
        expect(getExamNavigatorTargetIndex('Home', 2, 4)).toBe(0);
        expect(getExamNavigatorTargetIndex('End', 0, 4)).toBe(3);
        expect(getExamNavigatorTargetIndex('ArrowDown', 1, 4)).toBe(2);
        expect(getExamNavigatorTargetIndex('ArrowUp', 1, 4)).toBe(0);
        expect(getExamNavigatorTargetIndex('Escape', 1, 4)).toBe(-1);
    });
});

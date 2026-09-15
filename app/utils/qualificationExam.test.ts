import { describe, expect, it } from 'vitest';
import { isValidExamResponseKind } from './qualificationExam';

describe('isValidExamResponseKind', () => {
    it.each([
        ['freeText', 'freeText', true],
        ['singleChoice', 'freeText', false],
        ['image', 'separator', true],
        ['image', 'image', false],
        [undefined, 'separator', false],
        ['yesno', undefined, false],
    ])('validates %s -> %s', (questionKind, responseKind, expected) => {
        expect(isValidExamResponseKind(questionKind, responseKind)).toBe(expected);
    });
});

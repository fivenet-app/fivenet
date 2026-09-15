import { describe, expect, it } from 'vitest';
import { isAnsweredExamSingleChoice, isValidExamResponseKind } from './qualificationExam';

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

describe('isAnsweredExamSingleChoice', () => {
    it.each([
        ['A', true],
        ['', false],
        ['C', false],
        [false, false],
    ])('accepts configured string choices only', (choice, expected) => {
        expect(isAnsweredExamSingleChoice(choice, ['A', 'B'])).toBe(expected);
    });
});

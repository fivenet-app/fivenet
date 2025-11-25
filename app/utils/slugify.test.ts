import { describe, expect, it } from 'vitest';
import slug from './slugify';

describe('slugify', () => {
    it('should convert a string to a slug', () => {
        const input = 'Hello World!';
        const expected = 'hello-world';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with special characters', () => {
        const input = 'Café & Restaurant';
        const expected = 'cafe-and-restaurant';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with multiple spaces', () => {
        const input = '   Multiple   Spaces   ';
        const expected = 'multiple-spaces';
        expect(slug(input)).toBe(expected);
    });

    it('should handle empty strings', () => {
        const input = '';
        const expected = '';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with non-ASCII characters', () => {
        const input = '你好，世界';
        const expected = '';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with trailing dashes', () => {
        const input = 'Trailing Dash-';
        const expected = 'trailing-dash';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with leading dashes', () => {
        const input = '-Leading Dash';
        const expected = 'leading-dash';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with multiple dashes', () => {
        const input = 'Multiple---Dashes';
        const expected = 'multiple-dashes';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with punctuation marks', () => {
        const input = 'Multiple. Dashes?!? Here!';
        const expected = 'multiple-dashes-here';
        expect(slug(input)).toBe(expected);
    });

    it('should handle strings with ampersign and semi colons', () => {
        const input = 'Amper & Sign; Semi; Colons';
        const expected = 'amper-and-sign-semi-colons';
        expect(slug(input)).toBe(expected);
    });
});

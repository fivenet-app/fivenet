import { describe, expect, it } from 'vitest';
import { backgroundColors, hexToRgb, isColorBright, primaryColors, stringToColor } from './color';

describe('primaryColors', () => {
    it('should contain predefined primary colors', () => {
        expect(primaryColors).toContainEqual({
            label: 'green',
            chip: { color: 'green' },
            class: 'bg-green-500 dark:bg-green-400',
        });
    });
});

describe('backgroundColors', () => {
    it('should contain predefined background colors', () => {
        expect(backgroundColors).toContainEqual({
            label: 'slate',
            chip: { color: 'slate' },
            class: 'bg-slate-500 dark:bg-slate-400',
        });
    });
});

describe('stringToColor', () => {
    it('should generate a consistent color for a given string', () => {
        const str = 'test';
        const color = stringToColor(str);
        expect(color).toMatch(/^#[0-9a-fA-F]{6}$/);
    });
});

describe('hexToRgb', () => {
    it('should convert a hex color to RGB', () => {
        const hex = '#ff5733';
        const rgb = hexToRgb(hex);
        expect(rgb).toEqual({ r: 255, g: 87, b: 51 });
    });

    it('should return default value if hex is invalid', () => {
        const hex = 'invalid';
        const def = { r: 0, g: 0, b: 0 };
        const rgb = hexToRgb(hex, def);
        expect(rgb).toEqual(def);
    });
});

describe('isColorBright', () => {
    it('should return true for bright colors', () => {
        const brightColor = '#ffffff';
        expect(isColorBright(brightColor)).toBe(true);
    });

    it('should return false for dark colors', () => {
        const darkColor = '#000000';
        expect(isColorBright(darkColor)).toBe(false);
    });

    it('should handle RGB input', () => {
        const rgb = { r: 255, g: 255, b: 255 };
        expect(isColorBright(rgb)).toBe(true);
    });
});

export const primaryColors = [
    // Primary - Default
    { label: 'green', chip: { color: 'green' } },
    { label: 'teal', chip: { color: 'teal' } },
    { label: 'cyan', chip: { color: 'cyan' } },
    { label: 'sky', chip: { color: 'sky' } },
    { label: 'blue', chip: { color: 'blue' } },
    { label: 'indigo', chip: { color: 'indigo' } },
    { label: 'violet', chip: { color: 'violet' } },
    // Primary - Custom
    { label: 'yellow', chip: { color: 'yellow' } },
    { label: 'amber', chip: { color: 'amber' } },
    { label: 'lime', chip: { color: 'lime' } },
    { label: 'emerald', chip: { color: 'emerald' } },
    { label: 'fuchsia', chip: { color: 'fuchsia' } },
    { label: 'rose', chip: { color: 'rose' } },
    { label: 'pink', chip: { color: 'pink' } },
    { label: 'orange', chip: { color: 'orange' } },
    { label: 'red', chip: { color: 'red' } },
    { label: 'purple', chip: { color: 'purple' } },
] as const satisfies readonly Color[];

export const backgroundColors = [
    // Gray Colors
    { label: 'slate', chip: { color: 'slate' } },
    { label: 'zinc', chip: { color: 'zinc' } },
    { label: 'neutral', chip: { color: 'neutral' } },
    { label: 'stone', chip: { color: 'stone' } },
    { label: 'taupe', chip: { color: 'taupe' } },
    { label: 'mauve', chip: { color: 'mauve' } },
    { label: 'mist', chip: { color: 'mist' } },
    { label: 'olive', chip: { color: 'olive' } },
] as const satisfies readonly Color[];

export type PaletteColor = (typeof primaryColors)[number]['label'] | (typeof backgroundColors)[number]['label'];

const nuxtUiSemanticColors = ['primary', 'secondary', 'success', 'info', 'warning', 'error', 'gray'] as const;

export type NuxtUIColor = (typeof nuxtUiSemanticColors)[number] | PaletteColor;

export const nuxtUiColors = [
    ...nuxtUiSemanticColors,
    ...primaryColors.map((color) => color.label as NuxtUIColor),
    ...backgroundColors.map((color) => color.label as NuxtUIColor),
] as const satisfies readonly NuxtUIColor[];

export const rgbBlack = { r: 0, g: 0, b: 0 };

// Taken from https://stackoverflow.com/a/16348977
export function stringToColor(str: string): string {
    let hash = 0;
    str.split('').forEach((char) => {
        hash = char.charCodeAt(0) + ((hash << 5) - hash);
    });
    let color = '#';
    for (let i = 0; i < 3; i++) {
        const value = (hash >> (i * 8)) & 0xff;
        color += value.toString(16).padStart(2, '0');
    }
    return color;
}

// Taken from https://stackoverflow.com/a/5624139
export function hexToRgb(hex: string, def: RGB | undefined = undefined): RGB | undefined {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);

    return result
        ? {
              r: result[1] ? parseInt(result[1], 16) : 0,
              g: result[2] ? parseInt(result[2], 16) : 0,
              b: result[3] ? parseInt(result[3], 16) : 0,
          }
        : def;
}

export function isColorBright(input: RGB | string): boolean {
    const rgb = typeof input === 'string' ? hexToRgb(input, rgbBlack)! : input;

    // http://www.w3.org/TR/AERT#color-contrast
    const brightness = Math.round((rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000);
    return brightness > 125;
}

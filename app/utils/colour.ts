// Taken from https://stackoverflow.com/a/16348977
export function stringToColour(str: string): string {
    let hash = 0;
    str.split('').forEach((char) => {
        hash = char.charCodeAt(0) + ((hash << 5) - hash);
    });
    let colour = '#';
    for (let i = 0; i < 3; i++) {
        const value = (hash >> (i * 8)) & 0xff;
        colour += value.toString(16).padStart(2, '0');
    }
    return colour;
}

export type RGB = { r: number; g: number; b: number };
export const RGBBlack = { r: 0, g: 0, b: 0 };

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

export function isColourBright(input: RGB | string): boolean {
    const rgb = typeof input === 'string' ? hexToRgb(input, RGBBlack)! : input;

    // http://www.w3.org/TR/AERT#color-contrast
    const brightness = Math.round((rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000);
    return brightness > 125;
}

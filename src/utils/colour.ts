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

// Taken from https://stackoverflow.com/a/5624139
export function hexToRgb(hex: string): RGB | undefined {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result
        ? {
              r: parseInt(result[1], 16),
              g: parseInt(result[2], 16),
              b: parseInt(result[3], 16),
          }
        : undefined;
}

export function colourToTextColour(rgb: RGB): string {
    // http://www.w3.org/TR/AERT#color-contrast
    const brightness = Math.round((rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000);
    const textColour = brightness > 125 ? 'black' : 'yellow';
    return textColour;
}

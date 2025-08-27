export const primaryColors: Color[] = [
    // Primary - Default
    { label: 'green', chip: { color: 'green' }, class: 'bg-green-500 dark:bg-green-400' },
    { label: 'teal', chip: { color: 'teal' }, class: 'bg-teal-500 dark:bg-teal-400' },
    { label: 'cyan', chip: { color: 'cyan' }, class: 'bg-cyan-500 dark:bg-cyan-400' },
    { label: 'sky', chip: { color: 'sky' }, class: 'bg-sky-500 dark:bg-sky-400' },
    { label: 'blue', chip: { color: 'blue' }, class: 'bg-blue-500 dark:bg-blue-400' },
    { label: 'indigo', chip: { color: 'indigo' }, class: 'bg-indigo-500 dark:bg-indigo-400' },
    { label: 'violet', chip: { color: 'violet' }, class: 'bg-violet-500 dark:bg-violet-400' },
    // Primary - Custom
    { label: 'yellow', chip: { color: 'yellow' }, class: 'bg-yellow-500 dark:bg-yellow-400' },
    { label: 'amber', chip: { color: 'amber' }, class: 'bg-amber-500 dark:bg-amber-400' },
    { label: 'lime', chip: { color: 'lime' }, class: 'bg-lime-500 dark:bg-lime-400' },
    { label: 'emerald', chip: { color: 'emerald' }, class: 'bg-emerald-500 dark:bg-emerald-400' },
    { label: 'fuchsia', chip: { color: 'fuchsia' }, class: 'bg-fuchsia-500 dark:bg-fuchsia-400' },
    { label: 'rose', chip: { color: 'rose' }, class: 'bg-rose-500 dark:bg-rose-400' },
    { label: 'pink', chip: { color: 'pink' }, class: 'bg-pink-500 dark:bg-pink-400' },
    { label: 'orange', chip: { color: 'orange' }, class: 'bg-orange-500 dark:bg-orange-400' },
    { label: 'error', chip: { color: 'error' }, class: 'bg-red-500 dark:bg-red-400' },
    { label: 'purple', chip: { color: 'purple' }, class: 'bg-purple-500 dark:bg-purple-400' },
] as const;

export const backgroundColors: Color[] = [
    // Gray Colors
    { label: 'slate', chip: { color: 'slate' }, class: 'bg-slate-500 dark:bg-slate-400' },
    { label: 'zinc', chip: { color: 'zinc' }, class: 'bg-zinc-500 dark:bg-zinc-400' },
    { label: 'neutral', chip: { color: 'neutral' }, class: 'bg-neutral-500 dark:bg-neutral-400' },
    { label: 'stone', chip: { color: 'stone' }, class: 'bg-stone-500 dark:bg-stone-400' },
] as const;

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

export function isColorBright(input: RGB | string): boolean {
    const rgb = typeof input === 'string' ? hexToRgb(input, RGBBlack)! : input;

    // http://www.w3.org/TR/AERT#color-contrast
    const brightness = Math.round((rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000);
    return brightness > 125;
}

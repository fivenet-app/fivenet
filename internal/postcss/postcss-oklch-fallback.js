import { formatRgb, parse } from 'culori';

export default () => ({
    postcssPlugin: 'postcss-oklch-fallback',
    Declaration(decl) {
        if (decl.value.includes('oklch(')) {
            // Replace each occurrence of oklch(...) in the value
            const fallbackValue = decl.value.replace(/oklch\(([^)]+)\)/g, (_, contents) => {
                try {
                    // Parse the oklch color using culori. The string must be in a format culori understands.
                    const color = parse(`oklch(${contents})`);
                    // Convert it to an sRGB string (like 'rgb(…)')
                    const rgb = formatRgb(color);
                    return rgb;
                } catch (e) {
                    console.error('Error converting oklch color:', e);
                    // In case of error, leave the value as-is.
                    return `oklch(${contents})`;
                }
            });

            // Insert a fallback declaration BEFORE the current one.
            // Browsers that don’t understand oklch will use this fallback.
            decl.cloneBefore({ value: fallbackValue });
        }
    },
});

export const postcss = true;

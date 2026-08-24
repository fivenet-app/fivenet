const viewportUnitPattern = /\b(\d*\.?\d+)(svh|dvh|lvh)\b/g;

const viewportHeightProperties = new Set(['height', 'min-height', 'max-height']);

function replaceViewportHeightUnits(value) {
    const fallbackValue = value.replaceAll(viewportUnitPattern, '$1vh');
    return fallbackValue === value ? undefined : fallbackValue;
}

export default function viewportUnitFixup() {
    return {
        postcssPlugin: 'postcss-viewport-unit-fixup',
        Declaration(decl) {
            if (!viewportHeightProperties.has(decl.prop)) return;

            const fallbackValue = replaceViewportHeightUnits(decl.value);
            if (fallbackValue === undefined || fallbackValue === decl.value) return;

            decl.cloneBefore({ value: fallbackValue });
        },
    };
}

export const postcss = true;

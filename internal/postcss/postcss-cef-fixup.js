// Chromium Embedded Framework (CEF) in FiveM is too old to support the syntax without a space after the first value
export default function cefFixup() {
    return {
        postcssPlugin: 'postcss-cef-fixup',
        Declaration(decl) {
            if (decl.prop !== '--tw-ring-shadow') return;

            // Replace "var(--tw-ring-offset-width,)" with "var(--tw-ring-offset-width, 0px)"
            decl.value = decl.value.replaceAll(/var\(\s*--tw-ring-offset-width\s*,?\)/g, 'var(--tw-ring-offset-width, 0px)');
            decl.value = decl.value.replaceAll(/var\(\s*--tw-ring-inset\s*,?\)/g, 'var(--tw-ring-inset, ) ');
        },
    };
}

export const postcss = true;

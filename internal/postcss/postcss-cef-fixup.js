module.exports = () => {
    return {
        postcssPlugin: 'postcss-cef-fixup',
        Declaration(decl) {
            if (decl.prop !== '--tw-ring-shadow') return;

            // Replace "var(--tw-ring-offset-width,)" with "var(--tw-ring-offset-width, 0px)"
            decl.value = decl.value.replace(/var\(\s*--tw-ring-offset-width\s*,?\)/g, 'var(--tw-ring-offset-width, 0px)');
            decl.value = decl.value.replace(/var\(\s*--tw-ring-inset,\)/g, 'var(--tw-ring-inset, )');
        },
    };
};

module.exports.postcss = true;

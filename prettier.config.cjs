/** @type {import("prettier").Config} */
module.exports = {
    useTabs: false,
    tabWidth: 4,
    semi: true,
    singleQuote: true,
    printWidth: 128,
    plugins: ['prettier-plugin-tailwindcss', 'prettier-plugin-organize-imports'],
    tailwindConfig: './tailwind.config.ts',
};

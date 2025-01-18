/** @type {import("prettier").Config} */
module.exports = {
    useTabs: false,
    tabWidth: 4,
    semi: true,
    singleQuote: true,
    printWidth: 128,
    plugins: ['prettier-plugin-organize-imports', 'prettier-plugin-tailwindcss'],
    tailwindConfig: './tailwind.config.ts',
    tailwindStylesheet: './app/assets/css/tailwind.css',
};

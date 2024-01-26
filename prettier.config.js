/** @type {import("prettier").Config} */
const config = {
    useTabs: false,
    tabWidth: 4,
    semi: true,
    singleQuote: true,
    printWidth: 128,
    plugins: ['prettier-plugin-tailwindcss', 'prettier-plugin-organize-imports'],
    tailwindConfig: './tailwind.config.cjs',
};

export default config;

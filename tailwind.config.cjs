/** @type {import('tailwindcss').Config} */

const defaultTheme = require('tailwindcss/defaultTheme');

module.exports = {
    mode: 'jit',
    content: ['./public/**/*.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
    plugins: [require('@tailwindcss/forms'), require('@tailwindcss/typography')],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter var', ...defaultTheme.fontFamily.sans],
            },
            colors: {
                'dark-primary': '#0f0f0f',
            },
        },
    },
};

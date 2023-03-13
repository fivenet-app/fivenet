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
            boxShadow: {
                float: '0px 0px 22px 0px rgba(0,0,0,0.5)'
            },
            colors: {
                neutral: '#ffffff',
                base: {
                    100: '#e2e3e6',
                    200: '#c4c6cd',
                    300: '#a7aab4',
                    400: '#898d9b',
                    500: '#6c7182',
                    600: '#565a68',
                    700: '#41444e',
                    800: '#2b2d34',
                    850: '#202227',
                    900: '#16171a',
                },
            },
        },
    },
};

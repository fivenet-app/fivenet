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
                success: {
                    100: '#f4fff6',
                    200: '#eaffec',
                    300: '#dfffe3',
                    400: '#d5ffd9',
                    500: '#caffd0',
                    600: '#a2cca6',
                    700: '#79997d',
                    800: '#516653',
                    900: '#28332a',
                },
                info: {
                    100: '#dce5f2',
                    200: '#b9cbe4',
                    300: '#96b1d7',
                    400: '#7397c9',
                    500: '#507dbc',
                    600: '#406496',
                    700: '#304b71',
                    800: '#20324b',
                    900: '#101926',
                },
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

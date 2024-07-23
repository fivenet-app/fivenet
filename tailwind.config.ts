import type { Config } from 'tailwindcss';
import defaultTheme from 'tailwindcss/defaultTheme';

export default <Partial<Config>>{
    mode: 'jit',
    content: [
        `./app/components/**/*.{vue,js,ts}`,
        `./app/layouts/**/*.vue`,
        `./app/pages/**/*.vue`,
        `./app/composables/**/*.{js,ts}`,
        `./app/plugins/**/*.{js,ts}`,
        `./app/store/**/*.{js,ts}`,
        `./app/utils/**/*.{js,ts}`,
        `./app/App.{js,ts,vue}`,
        `./app/app.{js,ts,vue}`,
        `./app/Error.{js,ts,vue}`,
        `./app/error.{js,ts,vue}`,
        `./app/app.config.{js,ts}`,
        `./nuxt.config.{js,ts}`,
    ],
    theme: {
        extend: {
            colors: {
                base: {
                    100: '#e2e3e6',
                    200: '#c4c6cd',
                    300: '#a7aab4',
                    400: '#898d9b',
                    500: '#6c7182',
                    600: '#565a68',
                    700: '#41444e',
                    800: '#2b2d34',
                    900: '#16171a',
                },
                error: {
                    100: '#f7d4d7',
                    200: '#efa8af',
                    300: '#e77d88',
                    400: '#df5160',
                    500: '#d72638',
                    600: '#ac1e2d',
                    700: '#811722',
                    800: '#560f16',
                    900: '#2b080b',
                },
                warn: {
                    100: '#fdefd5',
                    200: '#fbdfab',
                    300: '#face81',
                    400: '#f8be57',
                    500: '#f6ae2d',
                    600: '#c58b24',
                    700: '#94681b',
                    800: '#624612',
                    900: '#312309',
                },
                info: {
                    100: '#ddf8f9',
                    200: '#bbf1f3',
                    300: '#99ebec',
                    400: '#77e4e6',
                    500: '#55dde0',
                    600: '#44b1b3',
                    700: '#338586',
                    800: '#22585a',
                    900: '#112c2d',
                },
                success: {
                    100: '#d3f5e1',
                    200: '#a7ebc4',
                    300: '#7be2a6',
                    400: '#4fd889',
                    500: '#23ce6b',
                    600: '#1ca556',
                    700: '#157c40',
                    800: '#0e522b',
                    900: '#072915',
                },
            },
            fontFamily: {
                sans: ['DM Sans', ...defaultTheme.fontFamily.sans],
            },
            boxShadow: {
                float: '0px 0px 22px 0px rgba(0,0,0,0.5)',
                glow: '0px 0px 22px 0px rgba(255,255,255,0.2)',
            },
            maxWidth: {
                '8xl': '88rem',
                screen: '100vw',
            },
            height: {
                '112': '28rem',
                dscreen: '100dvh',
            },
            animation: {
                wiggle: 'wiggle 1s infinite',
            },
            keyframes: {
                wiggle: {
                    '0%': {
                        transform: 'rotate(0deg)',
                    },
                    '80%': {
                        transform: 'rotate(0deg)',
                    },
                    '85%': {
                        transform: 'rotate(5deg)',
                    },
                    '95%': {
                        transform: 'rotate(-5deg)',
                    },
                    '100%': {
                        transform: 'rotate(0deg)',
                    },
                },
            },
            zIndex: {
                '150': '150',
                '100': '100',
                '5': '5',
            },
        },
    },
};

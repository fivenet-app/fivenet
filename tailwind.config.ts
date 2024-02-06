import type { Config } from 'tailwindcss';
import defaultTheme from 'tailwindcss/defaultTheme';

export default <Partial<Config>>{
    mode: 'jit',
    content: [
        `./src/components/**/*.{vue,js,ts}`,
        `./src/layouts/**/*.vue`,
        `./src/pages/**/*.vue`,
        `./src/composables/**/*.{js,ts}`,
        `./src/plugins/**/*.{js,ts}`,
        `./src/store/**/*.{js,ts}`,
        `./src/utils/**/*.{js,ts}`,
        `./src/App.{js,ts,vue}`,
        `./src/app.{js,ts,vue}`,
        `./src/Error.{js,ts,vue}`,
        `./src/error.{js,ts,vue}`,
        `./src/app.config.{js,ts}`,
        `./nuxt.config.{js,ts}`,
    ],
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
        // eslint-disable-next-line @typescript-eslint/no-var-requires
        require('tailwindcss-themer')({
            defaultTheme: {
                extend: {
                    colors: {
                        black: '#000000',
                        white: '#ffffff',
                        neutral: '#ffffff',
                        'body-color': '#16171a', // `bg-base-900` from `defaultTheme`
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
                        primary: {
                            100: '#e3dffc',
                            200: '#c6c0f9',
                            300: '#aaa0f5',
                            400: '#8d81f2',
                            500: '#7161ef',
                            600: '#5a4ebf',
                            700: '#443a8f',
                            800: '#2d2760',
                            900: '#171330',
                        },
                        secondary: {
                            100: '#dbd7e5',
                            200: '#b7afca',
                            300: '#9387b0',
                            400: '#6f5f95',
                            500: '#4b377b',
                            600: '#3c2c62',
                            700: '#2d214a',
                            800: '#1e1631',
                            900: '#0f0b19',
                        },
                        accent: {
                            100: '#d4d5e8',
                            200: '#a9abd0',
                            300: '#7d80b9',
                            400: '#5256a1',
                            500: '#272c8a',
                            600: '#1f236e',
                            700: '#171a53',
                            800: '#101237',
                            900: '#08091c',
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
                },
            },
            themes: [
                {
                    name: 'themePurple',
                    extend: {
                        colors: {
                            'body-color': '#101129',
                            primary: {
                                100: '#e3dffc',
                                200: '#c6c0f9',
                                300: '#aaa0f5',
                                400: '#8d81f2',
                                500: '#7161ef',
                                600: '#5a4ebf',
                                700: '#443a8f',
                                800: '#2d2760',
                                900: '#171330',
                            },
                            secondary: {
                                100: '#dbd7e5',
                                200: '#b7afca',
                                300: '#9387b0',
                                400: '#6f5f95',
                                500: '#4b377b',
                                600: '#3c2c62',
                                700: '#2d214a',
                                800: '#1e1631',
                                900: '#0f0b19',
                            },
                            accent: {
                                100: '#d4d5e8',
                                200: '#a9abd0',
                                300: '#7d80b9',
                                400: '#5256a1',
                                500: '#272c8a',
                                600: '#1f236e',
                                700: '#171a53',
                                800: '#101237',
                                900: '#08091c',
                            },
                        },
                    },
                },
                {
                    name: 'themeBaddieOrange',
                    extend: {
                        colors: {
                            'body-color': '#2a0f08',
                            // Based upon #0079a6
                            primary: {
                                100: '#cafbff',
                                200: '#9cf4ff',
                                300: '#57ebff',
                                400: '#0cd5ff',
                                500: '#00b9ea',
                                600: '#0091c4',
                                700: '#0079a6',
                                800: '#0b5d7f',
                                900: '#0e4d6b',
                            },
                            // Based upon #db561b
                            secondary: {
                                100: '#fcecd8',
                                200: '#f8d4b0',
                                300: '#f4b57d',
                                400: '#f0995c',
                                500: '#ea6e25',
                                600: '#db561b',
                                700: '#b64018',
                                800: '#91351b',
                                900: '#752c19',
                            },
                            // Based upon #9aa51d
                            accent: {
                                100: '#f8f8cf',
                                200: '#f0f1a5',
                                300: '#e2e571',
                                400: '#cfd645',
                                500: '#b3bc26',
                                600: '#9aa51d',
                                700: '#6a7219',
                                800: '#545b19',
                                900: '#474d1a',
                            },
                        },
                    },
                },
                {
                    name: 'themeBaddiePink',
                    extend: {
                        colors: {
                            'body-color': '#500724',
                            // Based upon #f472b6
                            primary: {
                                100: '#fce7f3',
                                200: '#fbcfe8',
                                300: '#f9a8d4',
                                400: '#f472b6',
                                500: '#ec4899',
                                600: '#db2777',
                                700: '#be185d',
                                800: '#9d174d',
                                900: '#831843',
                            },
                            // Based upon #f43f5e
                            secondary: {
                                100: '#fee5f6',
                                200: '#fecced',
                                300: '#ffa2df',
                                400: '#fe68c6',
                                500: '#f93bae',
                                600: '#e9198c',
                                700: '#cb0b70',
                                800: '#ab0d5e',
                                900: '#8b104e',
                            },
                            // Based upon #f9dc5c
                            accent: {
                                100: '#fdf7c4',
                                200: '#fbeb8d',
                                300: '#f9dc5c',
                                400: '#f4c31b',
                                500: '#e4ab0e',
                                600: '#c58309',
                                700: '#9d5d0b',
                                800: '#824a11',
                                900: '#6f3d14',
                            },
                        },
                    },
                },
                {
                    name: 'themeDaMedic',
                    extend: {
                        colors: {
                            'body-color': '#0d2744',
                            // Based upon #1e81ce
                            primary: {
                                100: '#e2effc',
                                200: '#bfddf8',
                                300: '#87c1f2',
                                400: '#47a2e9',
                                500: '#1e81ce',
                                600: '#126ab7',
                                700: '#105494',
                                800: '#11497b',
                                900: '#143e66',
                            },
                            // Based upon #9bc9ff
                            secondary: {
                                100: '#daeaff',
                                200: '#bedaff',
                                300: '#9bc9ff',
                                400: '#5da3fd',
                                500: '#377ffa',
                                600: '#2160ef',
                                700: '#194adc',
                                800: '#1b3db2',
                                900: '#1c388c',
                            },
                            // Based upon #a2f5c3
                            accent: {
                                100: '#d0fbdf',
                                200: '#a2f5c3',
                                300: '#6aeba6',
                                400: '#2fd882',
                                500: '#0bbe6a',
                                600: '#019a55',
                                700: '#017b48',
                                800: '#03623a',
                                900: '#045031',
                            },
                        },
                    },
                },
                {
                    name: 'themeBaddieYellow',
                    extend: {
                        colors: {
                            'body-color': '#54380a',
                            // Based upon #fde047
                            primary: {
                                100: '#fdffc1',
                                200: '#fffe86',
                                300: '#fff641',
                                400: '#ffe80d',
                                500: '#ffd900',
                                600: '#d1a000',
                                700: '#a67202',
                                800: '#89590a',
                                900: '#74490f',
                            },
                            // Based upon #f59e0b
                            secondary: {
                                100: '#feeac7',
                                200: '#fdd28a',
                                300: '#fcbb4d',
                                400: '#fbab24',
                                500: '#f59e0b',
                                600: '#d98b06',
                                700: '#b47409',
                                800: '#92610e',
                                900: '#78510f',
                            },
                            // Based upon #10b981
                            accent: {
                                100: '#d1faec',
                                200: '#a7f3da',
                                300: '#6ee7bf',
                                400: '#34d39e',
                                500: '#10b981',
                                600: '#059666',
                                700: '#047852',
                                800: '#065f42',
                                900: '#064e36',
                            },
                        },
                    },
                },
            ],
        }),
    ],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter', ...defaultTheme.fontFamily.sans],
            },
            boxShadow: {
                float: '0px 0px 22px 0px rgba(0,0,0,0.5)',
                glow: '0px 0px 22px 0px rgba(255,255,255,0.2)',
            },
            maxWidth: {
                '8xl': '88rem',
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

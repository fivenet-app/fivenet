import type { Config } from 'tailwindcss';
import defaultTheme from 'tailwindcss/defaultTheme';

export default <Partial<Config>>{
    theme: {
        extend: {
            fontFamily: {
                sans: ['DM Sans', ...defaultTheme.fontFamily.sans],
            },
            colors: {
                malibu: {
                    '50': '#f0f9ff',
                    '100': '#dff2ff',
                    '200': '#b8e5ff',
                    '300': '#65cbff',
                    '400': '#33bcfd',
                    '500': '#09a4ee',
                    '600': '#0082cc',
                    '700': '#0068a5',
                    '800': '#045888',
                    '900': '#0a4970',
                    '950': '#062e4b',
                },
            },
        },
    },
};

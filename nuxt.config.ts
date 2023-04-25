import { defineNuxtConfig } from 'nuxt/config';
import mkcert from 'vite-plugin-mkcert';
import fs from 'fs';
import path from 'path';

type PackageJson = {
    version: string;
};

const packageJson = fs.readFileSync('./package.json');
const version: string = (JSON.parse(packageJson.toString()) as PackageJson).version || '0.0.0';

mkcert();

// https://nuxt.com/docs/api/configuration/nuxt-config
const config = defineNuxtConfig({
    srcDir: 'src/',
    telemetry: false,
    ssr: false,
    modules: ['@pinia/nuxt', '@pinia-plugin-persistedstate/nuxt', 'nuxt-typed-router', '@nuxtjs/i18n', '@nuxtjs/tailwindcss'],
    pinia: {
        autoImports: [
            // automatically imports `defineStore`
            'defineStore', // import { defineStore } from 'pinia'
            ['defineStore', 'definePiniaStore'], // import { defineStore as definePiniaStore } from 'pinia'
            'acceptHMRUpdate', // import { acceptHMRUpdate } from 'pinia'
            'storeToRefs',
        ],
    },
    piniaPersistedstate: {
        storage: 'localStorage',
        debug: true,
    },
    i18n: {
        vueI18n: './i18n.config.ts',
        strategy: 'no_prefix',
        detectBrowserLanguage: {
            useCookie: false,
        },
        locales: [
            {
                name: 'English',
                code: 'en',
                iso: 'en',
                file: 'en-US.json',
            },
            {
                name: 'German',
                code: 'de',
                iso: 'de',
                file: 'de-DE.json',
            },
        ],
        lazy: false,
        langDir: './lang',
        defaultLocale: 'en',
    },
    vite: {
        define: {
            __APP_VERSION__: '"' + version + '"',
        },
        build: {
            commonjsOptions: {
                transformMixedEsModules: true,
            },
            manifest: true,
        },
        server: {
            hmr: {
                protocol: 'ws',
            },
            https: false,
            proxy: {
                '/api': 'http://localhost:8080',
                '/grpc': {
                    target: 'http://localhost:8181',
                },
            },
        },
    },
    css: [
        // Inter font (all weights)
        '@fontsource/inter/100.css',
        '@fontsource/inter/200.css',
        '@fontsource/inter/300.css',
        '@fontsource/inter/400.css',
        '@fontsource/inter/500.css',
        '@fontsource/inter/600.css',
        '@fontsource/inter/700.css',
        '@fontsource/inter/800.css',
        '@fontsource/inter/900.css',
    ],
    postcss: {
        plugins: {
            tailwindcss: {
                configPath: '~~/tailwind.config',
            },
            autoprefixer: {},
        },
    },
    typescript: {
        strict: true,
        tsConfig: {
            compilerOptions: {
                removeComments: true,
            },
        },
    },
    devServer: {},
    app: {
        head: {
            charset: 'utf-8',
            viewport: 'width=device-width, initial-scale=1',
            link: [{ rel: 'icon', type: 'image/png', href: '/images/logo.png' }],
        },
        pageTransition: { name: 'page', mode: 'out-in' },
    },
});

if (process.env.NODE_ENV !== 'production') {
    config.devServer!.https = {
        // Use vite-mkcert-plugin's cert + key for localhost
        key: path.resolve(process.env.HOME!, '.vite-plugin-mkcert', 'dev.pem'),
        cert: path.resolve(process.env.HOME!, '.vite-plugin-mkcert', 'cert.pem'),
    };
}

export default config;

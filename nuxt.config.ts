import fs from 'fs';
import { defineNuxtConfig } from 'nuxt/config';
import path from 'path';
import mkcert from 'vite-plugin-mkcert';
import { STRATEGIES } from 'vue-i18n-routing';

type PackageJson = {
    version: string;
};

const packageJson = fs.readFileSync('./package.json');
const version: string = (JSON.parse(packageJson.toString()) as PackageJson).version || '0.0.0';

const project = {
    name: 'FiveNet',
    shortName: 'FiveNet',
    description:
        "From searching the state's citizen and vehicles database, filling documents for investigations, court, and a livemap of your colleagues and dispatches. All that and more is (mostly) ready in this net, the FiveNet.",
    colors: {
        background: '#16171a',
        themeColor: '#1f236e',
    },
};

// https://nuxt.com/docs/api/configuration/nuxt-config
const config = defineNuxtConfig({
    srcDir: 'src/',
    telemetry: false,
    ssr: false,
    modules: [
        '@nuxt/devtools',
        '@pinia/nuxt',
        '@pinia-plugin-persistedstate/nuxt',
        'nuxt-typed-router',
        '@nuxtjs/robots',
        '@nuxtjs/i18n',
        '@nuxtjs/tailwindcss',
        '@vee-validate/nuxt',
        '@dargmuesli/nuxt-cookie-control',
    ],
    devtools: {
        enabled: true,
        vscode: {},
        timeline: {
            enabled: true,
        },
    },
    pinia: {
        autoImports: [
            // automatically imports `defineStore`
            'defineStore', // import { defineStore } from 'pinia'
            ['defineStore', 'definePiniaStore'], // import { defineStore as definePiniaStore } from 'pinia'
            'acceptHMRUpdate', // import { acceptHMRUpdate } from 'pinia'
            'storeToRefs',
        ],
    },
    robots: {
        rules: {
            UserAgent: '*',
            Disallow: '/',
            Allow: ['/$', '/index.html'],
        },
    },
    piniaPersistedstate: {
        storage: 'localStorage',
        debug: false,
    },
    i18n: {
        vueI18n: './i18n.config.ts',
        strategy: STRATEGIES.NO_PREFIX,
        detectBrowserLanguage: false,
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
        debug: false,
        lazy: true,
        langDir: './lang',
        defaultLocale: 'de',
        defaultDirection: 'ltr',
        baseUrl: '',
        trailingSlash: false,
        types: 'composition',
        compilation: {
            strictMessage: false,
            jit: true,
        },
        skipSettingLocaleOnNavigate: true,
        parallelPlugin: true,
    },
    veeValidate: {
        // disable or enable auto imports
        autoImports: true,
        // Use different names for components
        componentNames: {
            Form: 'VeeForm',
            Field: 'VeeField',
            FieldArray: 'VeeFieldArray',
            ErrorMessage: 'VeeErrorMessage',
        },
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
                protocol: 'wss',
            },
            https: true,
            proxy: {
                '/api': 'http://localhost:8080',
                '/grpc': {
                    target: 'http://localhost:8181',
                    // Make sure streaming works, but is also limited by the "real world" (3600s = 60m)
                    proxyTimeout: 3600 * 1000,
                    timeout: 3600 * 1000,
                },
            },
        },
        plugins: [mkcert()],
        optimizeDeps: {
            exclude: ['vue-demi'],
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
            'postcss-import': {},
            'tailwindcss/nesting': 'postcss-nesting',
            tailwindcss: {},
            autoprefixer: {},
            ...(process.env.NODE_ENV === 'production' ? { cssnano: {} } : {}),
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
    devServer: {
        port: 3000,
        https: {
            // Use vite-mkcert-plugin's cert + key for localhost
            key: path.resolve(process.env.HOME!, '.vite-plugin-mkcert', 'dev.pem'),
            cert: path.resolve(process.env.HOME!, '.vite-plugin-mkcert', 'cert.pem'),
        },
    },
    app: {
        head: {
            charset: 'utf-8',
            viewport: 'width=device-width, initial-scale=1',
            link: [{ rel: 'icon', type: 'image/png', href: '/images/logo.png' }],
            htmlAttrs: {
                class: 'h-full',
            },
            bodyAttrs: {
                class: 'h-full',
            },
        },
        pageTransition: { name: 'page', mode: 'out-in' },
    },
    cookieControl: {
        barPosition: 'bottom-full',
        closeModalOnClickOutside: true,
        isControlButtonEnabled: true,
        colors: {
            // Tailwind CSS - bg-accent-600
            modalBackground: '#1f236e',
            controlButtonBackground: '#1f236e',
            // Tailwind CSS - bg-accent-600 + bg-accent-100/10
            controlButtonHoverBackground: '#32357b',
            controlButtonIconColor: '#fff',
            controlButtonIconHoverColor: '#fff',
            modalTextColor: '#fff',
        },
        locales: ['en', 'de'],
        cookies: {
            necessary: [],
            optional: [
                {
                    name: 'Social Login Cookies',
                    description: 'Cookies used for FiveNet Social Login functionality.',
                    id: 'social_login',
                    targetCookieIds: ['fivenet_oauth2_state', 'fivenet_token'],
                },
            ],
        },
    },
});

export default config;

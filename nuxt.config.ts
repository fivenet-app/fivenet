import { STRATEGIES } from 'vue-i18n-routing';

const appVersion: string = process.env.COMMIT_REF || 'COMMIT_REF';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    telemetry: false,
    ssr: false,
    extends: [process.env.NUXT_UI_PRO_PATH || '@nuxt/ui-pro'],

    modules: [
        '@nuxt/ui',
        'nuxt-typed-router',
        '@nuxtjs/robots',
        '@nuxt/fonts',
        '@vueuse/nuxt',
        '@pinia/nuxt',
        'pinia-plugin-persistedstate/nuxt',
        'nuxt-zod-i18n',
        '@nuxtjs/i18n',
        '@galexrt/nuxt-update',
        '@nuxt/eslint',
        '@nuxt/content',
        'nuxt-tiptap-editor',
        '@vueuse/sound/nuxt',
    ],

    future: {
        compatibilityVersion: 4,
    },

    experimental: {
        defaults: {
            useAsyncData: {
                deep: true,
            },
        },
    },

    site: {
        indexable: false,
    },

    ui: {
        safelistColors: [
            // Primary - Default
            'primary',
            'green',
            'teal',
            'cyan',
            'sky',
            'blue',
            'indigo',
            'violet',
            // Custom
            'error',
            'warn',
            'info',
            'success',
            // Gray Colors
            'slate',
            'cool',
            'zinc',
            'neutral',
            'stone',
        ],
    },

    postcss: {
        plugins: {
            './internal/postcss/postcss-oklch-fallback': {},
        },
    },

    uiPro: {
        routerOptions: false,
    },

    icon: {
        collections: ['simple-icons', 'mdi', 'flagpack'],
        provider: 'iconify',
        iconifyApiEndpoint: '/api/icons',
        fallbackToApi: false,
        clientBundle: {
            scan: true,
        },
    },

    fonts: {
        families: [{ name: 'DM Sans', weights: [100, 200, 300, 400, 500, 600, 700, 800, 900], global: true }],
    },

    app: {
        head: {
            charset: 'utf-8',
            viewport: 'width=device-width, initial-scale=1',
            link: [{ rel: 'icon', type: 'image/png', href: '/images/logo.png' }],
            meta: [{ name: 'darkreader-lock', content: '' }],
        },
        pageTransition: { name: 'page', mode: 'out-in' },
    },

    colorMode: {
        preference: 'dark',
    },

    typescript: {
        strict: true,
        tsConfig: {
            compilerOptions: {
                removeComments: true,
            },
        },
    },

    imports: {
        presets: [
            {
                from: '@protobuf-ts/runtime-rpc',
                imports: [
                    {
                        name: 'RpcError',
                        type: true,
                    },
                ],
            },
        ],
    },

    sourcemap: {
        client: true,
    },

    devtools: {
        enabled: true,
    },

    vite: {
        define: {
            APP_VERSION: `"${appVersion}"`,
        },
        build: {
            target: 'es2016',
            terserOptions: {
                compress: {
                    // Only drop certain console log calls
                    drop_console: ['trace', 'debug', 'info'],
                    drop_debugger: true,
                },
            },
        },
        server: {
            proxy: {
                '/api/icons': {
                    target: 'https://api.iconify.design',
                    rewrite: (path: string) => path.replace(/^\/api\/icons/, ''),
                    changeOrigin: true,
                },
                '/api/grpc': {
                    target: 'http://localhost:8080',
                    ws: true,
                    proxyTimeout: 60 * 60 * 1000,
                    timeout: 60 * 60 * 1000,
                    changeOrigin: true,
                    configure: (proxy, options) => {
                        proxy.on('proxyReq', (_, req) => {
                            req.headers.origin =
                                typeof options.target === 'string'
                                    ? options.target
                                    : `${options.target!.protocol}//${options.target!.host}`;
                        });
                    },
                },
                '/api': {
                    target: 'http://localhost:8080',
                    changeOrigin: true,
                    configure: (proxy, options) => {
                        proxy.on('proxyReq', (_, req) => {
                            req.headers.origin =
                                typeof options.target === 'string'
                                    ? options.target
                                    : `${options.target!.protocol}//${options.target!.host}`;
                        });
                    },
                },
            },
        },
    },

    i18n: {
        lazy: true,
        parallelPlugin: true,
        strategy: STRATEGIES.NO_PREFIX,
        detectBrowserLanguage: false,
        defaultLocale: 'en',
        locales: [
            {
                code: 'en',
                language: 'en',
                name: 'English',
                isCatchallLocale: true,
                files: ['en/en.json', 'en/zod.json'],
                icon: 'i-flagpack-gb-ukm',
            },
            {
                name: 'German',
                code: 'de',
                language: 'de',
                files: ['de/de.json', 'de/zod.json'],
                icon: 'i-flagpack-de',
            },
        ],
        baseUrl: '',
        trailingSlash: false,
        compilation: {
            strictMessage: false,
        },
    },

    zodI18n: {
        useModuleLocale: true,
    },

    piniaPluginPersistedstate: {
        storage: 'localStorage',
    },

    update: {
        version: appVersion,
        checkInterval: 115,
        path: '/api/version',
    },

    tiptap: {},

    sound: {
        sounds: {
            scan: true,
        },
    },

    $production: {
        icon: {
            iconifyApiEndpoint: '/api/icons',
        },
    },

    $development: {
        icon: {
            iconifyApiEndpoint: 'https://api.iconify.design',
            clientBundle: {
                scan: false,
            },
        },
    },

    compatibilityDate: '2024-07-05',
});

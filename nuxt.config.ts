const appVersion: string = process.env.COMMIT_REF || 'COMMIT_REF';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    telemetry: false,
    ssr: false,

    modules: [
        '@nuxt/ui-pro',
        '@pinia/nuxt',
        'pinia-plugin-persistedstate/nuxt',
        'nuxt-typed-router',
        '@nuxtjs/robots',
        '@nuxt/fonts',
        '@nuxt/image',
        'nuxt-zod-i18n',
        '@nuxtjs/i18n',
        '@galexrt/nuxt-update',
        '@nuxt/eslint',
        '@nuxt/content',
        'nuxt-tiptap-editor',
        '@vueuse/nuxt',
        '@vueuse/sound/nuxt',
        '@nuxtjs/leaflet',
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
        spaLoadingTemplateLocation: 'body',

        granularCachedData: false,
    },

    site: {
        indexable: false,
    },

    css: ['~/assets/css/main.css'],

    ui: {
        theme: {
            colors: [
                // Theme colors
                'primary',
                'secondary',
                'success',
                'info',
                'warning',
                'error',
                // Palette colors
                'amber',
                'blue',
                'cyan',
                'emerald',
                'fuchsia',
                'green',
                'indigo',
                'lime',
                'orange',
                'pink',
                'purple',
                'red',
                'rose',
                'sky',
                'teal',
                'violet',
                'white',
                'yellow',
                // Gray Colors
                'gray',
                'neutral',
                'slate',
                'stone',
                'zinc',
            ],
        },
    },

    postcss: {
        plugins: {
            './internal/postcss/postcss-oklch-fallback': {},
        },
    },

    uiPro: {
        content: true,
    },

    image: {
        provider: 'none',
    },

    icon: {
        collections: ['simple-icons', 'lucide', 'mdi', 'flagpack'],
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
                    drop_console: false,
                    drop_debugger: true,
                },
            },
        },
        esbuild: {
            // Only drop debugger statements, console logs are kept
            drop: ['debugger'],
        },
        server: {
            // Make it easier to test the app behind a proxy server (e.g., ngrok), SSR is disabled
            // so we can use `allowedHosts: true` to allow any host
            allowedHosts: true,
            proxy: {
                '/api/grpc': {
                    target: 'http://127.0.0.1:8080',
                    changeOrigin: true,
                    configure: (proxy, options) => {
                        // changeOrigin doesn't work "correctly"
                        proxy.on('proxyReq', (_, req) => {
                            req.headers.origin =
                                typeof options.target === 'string'
                                    ? options.target
                                    : `${options.target!.protocol}//${options.target!.host}`;
                        });
                    },
                },
                '/api/grpcws': {
                    target: 'http://127.0.0.1:8080',
                    ws: true,
                    proxyTimeout: 60 * 60 * 1000,
                    timeout: 60 * 60 * 1000,
                    changeOrigin: true,
                    configure: (proxy, options) => {
                        // changeOrigin doesn't work "correctly"
                        proxy.on('proxyReq', (_, req) => {
                            req.headers.origin =
                                typeof options.target === 'string'
                                    ? options.target
                                    : `${options.target!.protocol}//${options.target!.host}`;
                        });
                    },
                },
                '/api': {
                    target: 'http://127.0.0.1:8080',
                    changeOrigin: true,
                    configure: (proxy, options) => {
                        // changeOrigin doesn't work "correctly"
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
        strategy: 'no_prefix',
        defaultLocale: 'en',
        detectBrowserLanguage: false,
        locales: [
            {
                code: 'en',
                language: 'en',
                name: 'English',
                isCatchallLocale: true,
                file: 'en.json',
                icon: 'i-flagpack-gb-ukm',
            },
            {
                name: 'German',
                code: 'de',
                language: 'de',
                file: 'de.json',
                icon: 'i-flagpack-de',
            },
        ],
        baseUrl: '',
        parallelPlugin: true,
        compilation: {
            strictMessage: false,
        },
    },

    zodI18n: {
        useModuleLocale: false,
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

    leaflet: {
        heat: true,
    },

    $production: {
        icon: {
            iconifyApiEndpoint: '/api/icons',
        },
    },

    $development: {
        icon: {
            iconifyApiEndpoint: 'https://api.iconify.design',
            provider: 'iconify',
            clientBundle: {
                scan: true,
            },
        },
    },

    compatibilityDate: '2024-07-05',
});

import { defineNuxtConfig } from 'nuxt/config';

const appVersion: string = process.env.COMMIT_REF || 'COMMIT_REF';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    telemetry: false,
    ssr: false,

    modules: [
        '@nuxt/ui',
        '@pinia/nuxt',
        'pinia-plugin-persistedstate/nuxt',
        'nuxt-typed-router',
        '@nuxtjs/robots',
        '@nuxt/fonts',
        '@nuxt/image',
        '@nuxtjs/i18n',
        '@galexrt/nuxt-update',
        '@nuxt/eslint',
        '@nuxt/content',
        'nuxt-tiptap-editor',
        '@vueuse/nuxt',
        '@nuxtjs/leaflet',
        '@nuxt/test-utils/module',
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
            './internal/postcss/postcss-color-mix-transparency-fallback': {},
            'postcss-preset-env': {
                stage: 2,
                features: {
                    'oklab-function': {
                        preserve: true,
                        enableProgressiveCustomProperties: true,
                        subFeatures: {
                            displayP3: false,
                        },
                    },
                    // Not in use by Nuxt UI (yet)
                    'random-function': false,
                    'sign-functions': false,
                    'stepped-value-functions': false,
                    'trigonometric-functions': false,
                },
                enableClientSidePolyfills: false,
                preserve: true,
                browsers: 'chrome >= 103',
            },
            './internal/postcss/postcss-cef-fixup': {},
        },
    },

    image: {
        provider: 'none',
    },

    icon: {
        collections: ['mdi', 'simple-icons', 'flagpack'],
        provider: 'iconify',
        iconifyApiEndpoint: '/api/icons',
        fallbackToApi: false,
        clientBundle: {
            scan: true,
        },
    },

    app: {
        head: {
            charset: 'utf-8',
            viewport: 'width=device-width, initial-scale=1',
            link: [
                { rel: 'shortcut icon', href: '/favicon.ico' },
                { rel: 'icon', type: 'image/png', href: '/images/logo.png' },
                { rel: 'apple-touch-icon', href: '/images/logo.png' },
            ],
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
        optimizeDeps: {
            include: [
                '@nuxt/ui > prosemirror-state',
                '@nuxt/ui > prosemirror-transform',
                '@nuxt/ui > prosemirror-model',
                '@nuxt/ui > prosemirror-view',
                '@nuxt/ui > prosemirror-gapcursor',
                '@zxcvbn-ts/core',
                '@vue/devtools-core',
                '@vue/devtools-kit',
                'css-blank-pseudo/browser',
                'css-has-pseudo/browser',
                'zod/v4',
                'uuid',
                'emoji-blast',
                'mitt',
                'date-fns',
                '@protobuf-ts/grpcweb-transport',
                '@protobuf-ts/runtime-rpc',
                '@protobuf-ts/runtime',
                'howler', // CJS
                'browser-headers', // CJS
                '@internationalized/date',
                'vue-json-pretty',
                'zod',
                'vue-timeline-chart',
                '@unovis/ts',
                '@unovis/vue',
                '@tanstack/vue-table',
                'vue-draggable-plus',
                'maska/vue',
                '@selemondev/vue3-signature-pad',
                'yjs',
                'diff-match-patch', // CJS
                '@tiptap/extension-collaboration',
                '@tiptap/extension-collaboration-caret',
                '@tiptap/extension-emoji',
                '@tiptap/extensions',
                '@tiptap/y-tiptap',
                'lib0/observable',
                'y-protocols/awareness',
                '@tiptap/extension-details',
                '@tiptap/extension-highlight',
                '@tiptap/extension-invisible-characters',
                '@tiptap/extension-list',
                '@tiptap/extension-node-range',
                '@tiptap/extension-subscript',
                '@tiptap/extension-superscript',
                '@tiptap/extension-table',
                '@tiptap/extension-text-align',
                '@tiptap/extension-text-style',
                '@tiptap/extension-underline',
                '@tiptap/extension-unique-id',
                'pica', // CJS
                '@tiptap/pm/state',
                '@tiptap/pm/view',
                '@samemichaeltadele/tiptap-compare',
                'vue-countup-v3',
                'mdi-vue3',
                'v-calendar',
            ],
        },
        server: {
            // Make it easier to test the app behind a proxy server (e.g., ngrok), SSR is disabled
            // so we can use `allowedHosts: true` to allow any host
            allowedHosts: true,
            proxy: {
                '/api/grpcws': {
                    target: 'http://localhost:8080',
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
                '/api/grpc': {
                    target: 'http://localhost:8080',
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
                    target: 'http://localhost:8080',
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
        parallelPlugin: true,
        compilation: {
            strictMessage: false,
        },

        locales: [
            {
                name: 'English',
                code: 'en',
                language: 'en',
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

    leaflet: {
        heat: true,
    },

    $development: {
        icon: {
            iconifyApiEndpoint: 'https://api.iconify.design',
        },
    },

    compatibilityDate: '2025-12-31',
});

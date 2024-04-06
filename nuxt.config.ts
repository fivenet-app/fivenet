import svgLoader from 'vite-svg-loader';
import { STRATEGIES } from 'vue-i18n-routing';

const appVersion: string = process.env.COMMIT_REF || 'COMMIT_REF';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    srcDir: 'src/',
    telemetry: false,
    ssr: false,
    extends: ['@nuxt/ui-pro'],
    // @ts-ignore ts2589 one of the modules probably has a typing issue
    modules: [
        'nuxt-typed-router',
        '@nuxtjs/robots',
        '@nuxt/ui',
        '@nuxtjs/tailwindcss',
        '@nuxtjs/i18n',
        '@pinia/nuxt',
        '@pinia-plugin-persistedstate/nuxt',
        '@vee-validate/nuxt',
        '@dargmuesli/nuxt-cookie-control',
        'nuxt-update',
        '@vueuse/nuxt',
    ],
    ui: {
        icons: ['heroicons', 'simple-icons', 'mdi', 'flagpack'],
        safelistColors: ['primary', 'gray', 'red', 'orange', 'green', 'error', 'warn', 'info', 'success'],
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
    css: [
        // DM Sans font (all weights)
        '@fontsource/dm-sans/100.css',
        '@fontsource/dm-sans/200.css',
        '@fontsource/dm-sans/300.css',
        '@fontsource/dm-sans/400.css',
        '@fontsource/dm-sans/500.css',
        '@fontsource/dm-sans/600.css',
        '@fontsource/dm-sans/700.css',
        '@fontsource/dm-sans/800.css',
        '@fontsource/dm-sans/900.css',
    ],
    tailwindcss: {
        exposeConfig: true,
        configPath: './tailwind.config.ts',
    },
    typescript: {
        typeCheck: false,
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
        vscode: {},
        timeline: {
            enabled: true,
        },
    },
    vite: {
        define: {
            __APP_VERSION__: `"${appVersion}"`,
        },
        build: {
            commonjsOptions: {
                transformMixedEsModules: true,
            },
            manifest: true,
            terserOptions: {
                compress: {
                    drop_console: true,
                    drop_debugger: true,
                },
            },
        },
        server: {
            hmr: {
                protocol: 'ws',
            },
            proxy: {
                '/api': 'http://localhost:8080',
                '/grpc': {
                    target: 'http://localhost:8181',
                    // Make sure streaming works, but is also limited by the "real world" (1800s = 30m)
                    proxyTimeout: 1800 * 1000,
                    timeout: 1800 * 1000,
                },
            },
        },
        optimizeDeps: {
            exclude: ['vue-demi'],
        },
        plugins: [svgLoader()],
    },
    robots: {
        rules: {
            UserAgent: '*',
            Disallow: '/',
            Allow: ['/$', '/index.html'],
        },
    },
    i18n: {
        vueI18n: './i18n.config.ts',
        strategy: STRATEGIES.NO_PREFIX,
        detectBrowserLanguage: false,
        skipSettingLocaleOnNavigate: true,
        locales: [
            {
                name: 'English',
                dir: 'ltr',
                isCatchallLocale: true,
                code: 'en',
                iso: 'en',
                files: ['en-US.json'],
                icon: 'i-flagpack-gb-ukm',
            },
            {
                name: 'German',
                code: 'de',
                iso: 'de',
                files: ['de-DE.json'],
                icon: 'i-flagpack-de',
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
        parallelPlugin: true,
    },
    piniaPersistedstate: {
        storage: 'localStorage',
        debug: false,
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
                    targetCookieIds: ['fivenet_oauth2_state'],
                    links: {
                        '/privacy_policy': 'Privacy Policy',
                        '/imprint': 'Imprint',
                    },
                },
            ],
        },
    },
    update: {
        version: appVersion,
        checkInterval: 110,
        path: '/api/version',
    },
});

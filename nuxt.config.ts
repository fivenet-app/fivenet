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
        'nuxt-zod-i18n',
        '@nuxtjs/i18n',
        '@pinia/nuxt',
        '@pinia-plugin-persistedstate/nuxt',
        '@vueuse/nuxt',
        'nuxt-update',
    ],
    ui: {
        icons: ['simple-icons', 'mdi', 'flagpack'],
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
    colorMode: {
        preference: 'dark',
    },
    typescript: {
        strict: true,
        tsConfig: {
            compilerOptions: {
                removeComments: true,
            },
            exclude: ['../docs'],
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
        vscode: {
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
            terserOptions: {
                compress: {
                    drop_console: true,
                    drop_debugger: true,
                },
            },
        },
        server: {
            proxy: {
                '/api/icons': {
                    target: 'https://api.iconify.design',
                    rewrite: (path) => path.replace(/^\/api\/icons/, ''),
                    changeOrigin: true,
                },
                '/api': 'http://localhost:8080',
                '/grpc': {
                    target: 'http://localhost:8181',
                    // Make sure streaming works, but is also limited by the "real world" (1800s = 30m)
                    proxyTimeout: 1800 * 1000,
                    timeout: 1800 * 1000,
                },
            },
        },
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
                files: ['en.json', 'en/zod.json'],
                icon: 'i-flagpack-gb-ukm',
            },
            {
                name: 'German',
                code: 'de',
                iso: 'de',
                files: ['de.json', 'de/zod.json'],
                icon: 'i-flagpack-de',
            },
        ],
        lazy: true,
        langDir: './lang',
        defaultLocale: 'de',
        defaultDirection: 'ltr',
        baseUrl: '',
        trailingSlash: false,
        compilation: {
            strictMessage: false,
        },
        parallelPlugin: true,
    },
    zodI18n: {
        useModuleLocale: true,
    },
    piniaPersistedstate: {
        storage: 'localStorage',
    },
    update: {
        version: appVersion,
        checkInterval: 110,
        path: '/api/version',
    },
});

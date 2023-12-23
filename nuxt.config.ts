import { STRATEGIES } from 'vue-i18n-routing';
import svgLoader from 'vite-svg-loader';

const commit: string = process.env.COMMIT_REF || 'COMMIT_REF';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    srcDir: 'src/',
    telemetry: false,
    ssr: false,
    modules: [
        '@nuxt/devtools',
        'nuxt-typed-router',
        '@nuxtjs/eslint-module',
        '@nuxtjs/robots',
        '@nuxtjs/tailwindcss',
        '@nuxtjs/i18n',
        '@pinia/nuxt',
        '@pinia-plugin-persistedstate/nuxt',
        '@vee-validate/nuxt',
        '@dargmuesli/nuxt-cookie-control',
        'nuxt-update',
    ],
    typescript: {
        typeCheck: true,
        strict: true,
        tsConfig: {
            compilerOptions: {
                removeComments: true,
            },
        },
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
    eslint: {
        lintOnStart: false,
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
        skipSettingLocaleOnNavigate: true,
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
            __APP_VERSION__: `"${commit}"`,
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
                    targetCookieIds: ['fivenet_oauth2_state'],
                },
            ],
        },
    },
    update: {
        version: commit,
        checkInterval: 90,
        path: '/api/version',
    },
    experimental: {
        emitRouteChunkError: 'automatic',
    },
});

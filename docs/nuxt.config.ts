import { STRATEGIES } from 'vue-i18n-routing';

const appVersion: string = process.env.COMMIT_REF || 'COMMIT_REF';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    telemetry: false,
    extends: ['@nuxt/ui-pro'],
    modules: [
        '@nuxt/content',
        '@nuxt/ui',
        '@nuxtjs/i18n',
        '@pinia/nuxt',
        '@pinia-plugin-persistedstate/nuxt',
        '@nuxthq/studio',
        'nuxt-og-image',
    ],
    hooks: {
        // Define `@nuxt/ui` components as global to use them in `.md` (feel free to add those you need)
        'components:extend': (components) => {
            const globals = components.filter((c) => ['UButton', 'UIcon'].includes(c.pascalName));

            globals.forEach((c) => (c.global = true));
        },
    },
    ui: {
        icons: ['mdi', 'simple-icons', 'flagpack'],
        safelistColors: ['primary', 'malibu'],
    },
    routeRules: {
        '/api/search.json': { prerender: true },
    },
    devtools: {
        enabled: true,
    },
    typescript: {
        strict: false,
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
    vite: {
        define: {
            __APP_VERSION__: `"${appVersion}"`,
        },
    },
    i18n: {
        strategy: STRATEGIES.PREFIX_EXCEPT_DEFAULT,
        detectBrowserLanguage: {
            useCookie: false,
        },
        skipSettingLocaleOnNavigate: true,
        locales: [
            {
                name: 'English',
                dir: 'ltr',
                isCatchallLocale: true,
                code: 'en',
                iso: 'en',
                files: ['en.json'],
                icon: 'i-flagpack-gb-ukm',
            },
            {
                name: 'German',
                code: 'de',
                iso: 'de',
                files: ['de.json'],
                icon: 'i-flagpack-de',
            },
        ],
        debug: false,
        lazy: true,
        langDir: './lang',
        defaultLocale: 'en',
        defaultDirection: 'ltr',
        baseUrl: '',
        trailingSlash: false,
        compilation: {
            strictMessage: false,
        },
        parallelPlugin: true,
    },
    piniaPersistedstate: {
        storage: 'localStorage',
        debug: false,
    },
});

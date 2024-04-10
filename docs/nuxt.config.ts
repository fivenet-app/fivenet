import { resolve } from 'path';
import { STRATEGIES } from 'vue-i18n-routing';

const appVersion: string = process.env.COMMIT_REF || 'COMMIT_REF';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    extends: ['@nuxt/ui-pro'],
    modules: [
        '@nuxt/content',
        '@nuxt/ui',
        '@nuxt/fonts',
        '@nuxtjs/i18n',
        '@pinia/nuxt',
        '@pinia-plugin-persistedstate/nuxt',
        '@nuxthq/studio',
        'nuxt-og-image',
    ],
    content: {
        sources: {
            content: {
                driver: 'fs',
                prefix: '/docs',
                base: resolve(__dirname, 'content'),
            },
        },
    },
    hooks: {
        // Define `@nuxt/ui` components as global to use them in `.md` (feel free to add those you need)
        'components:extend': (components) => {
            const globals = components.filter((c) => ['UButton', 'UIcon'].includes(c.pascalName));

            globals.forEach((c) => (c.global = true));
        },
    },
    ui: {
        icons: ['mdi', 'heroicons', 'simple-icons', 'flagpack'],
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

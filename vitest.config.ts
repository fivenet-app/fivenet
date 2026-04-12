import { defineVitestProject } from '@nuxt/test-utils/config';
import { configDefaults, defineConfig } from 'vitest/config';

export default defineConfig({
    root: '.',
    appType: 'spa',

    test: {
        // Stop Vue 3 from logging Suspense warnings in the console during tests
        onConsoleLog: (l) => {
            return !l.startsWith('<Suspense>');
        },

        projects: [
            await defineVitestProject({
                test: {
                    name: 'unit',
                    include: ['app/**/*.{test,spec}.ts'],
                    exclude: ['app/test/**', ...configDefaults.exclude, '.direnv/*'],
                    environment: 'nuxt',
                    hookTimeout: 30_000,
                },
            }),
            await defineVitestProject({
                test: {
                    name: 'nuxt',
                    include: ['app/test/nuxt/**/*.{test,spec}.ts'],
                    exclude: ['app/**/*.{e2e,unit}.{test,spec}.ts', ...configDefaults.exclude, '.direnv/*'],
                    environment: 'nuxt',
                    hookTimeout: 30_000,
                },
            }),
        ],
    },
});

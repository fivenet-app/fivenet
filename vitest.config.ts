import { defineVitestProject } from '@nuxt/test-utils/config';
import { configDefaults, defineConfig } from 'vitest/config';

export default defineConfig({
    root: '.',

    test: {
        projects: [
            await defineVitestProject({
                test: {
                    name: 'unit',
                    include: ['app/**/*.{test,spec}.ts'],
                    exclude: ['app/test/**', ...configDefaults.exclude, '.direnv/*'],
                    environment: 'nuxt',
                },
            }),
            await defineVitestProject({
                test: {
                    name: 'nuxt',
                    include: ['app/test/nuxt/**/*.{test,spec}.ts'],
                    exclude: ['app/**/*.{e2e,unit}.{test,spec}.ts', ...configDefaults.exclude, '.direnv/*'],
                    environment: 'nuxt',
                },
            }),
        ],
    },
});

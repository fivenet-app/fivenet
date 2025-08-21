import { configDefaults, defineConfig } from 'vitest/config';

export default defineConfig({
    test: {
        include: ['app/**/*.spec.ts'],
        exclude: [...configDefaults.exclude, '.direnv/*'],
    },
});

// @ts-check
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
import withNuxt from './.nuxt/eslint.config.mjs';

export default withNuxt(
    {
        ignores: ['gen/', 'proto/'],
    },
    eslintPluginPrettierRecommended,
    {
        rules: {
            'no-console': 0,
            'require-await': 0,
            '@typescript-eslint/no-unused-vars': [
                'warn',
                {
                    argsIgnorePattern: '^_',
                    varsIgnorePattern: '^_',
                    caughtErrorsIgnorePattern: '^_',
                },
            ],
            'vue/no-unused-vars': [
                'warn',
                {
                    ignorePattern: '^_',
                },
            ],
            '@typescript-eslint/unified-signatures': 'off',
            '@typescript-eslint/no-unused-expressions': 'off',
            'import/consistent-type-specifier-style': 'off',
            'import/no-duplicates': [
                'error',
                {
                    'prefer-inline': true,
                },
            ],
        },
    },
    {
        files: ['app/pages/**', 'app/layouts/**'],
        rules: {
            'vue/multi-word-component-names': 'off',
        },
    },
);

// @ts-check
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
// @ts-ignore no types available yet
import pluginVue from 'eslint-plugin-vue';
import withNuxt from './.nuxt/eslint.config.mjs';

export default withNuxt(
    {
        ignores: ['gen/', 'proto/'],
    },
    ...pluginVue.configs['flat/recommended'],
    eslintPluginPrettierRecommended,
    {
        rules: {
            'no-console': 0,
            'require-await': 0,
            'no-restricted-syntax': ['error', 'IfStatement > ExpressionStatement > AssignmentExpression'],
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
        },
    },
    {
        files: ['app/pages/**', 'app/layouts/**'],
        rules: {
            'vue/multi-word-component-names': 'off',
        },
    },
);

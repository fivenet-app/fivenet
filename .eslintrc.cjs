// .eslintrc.cjs
module.exports = {
    root: true,
    env: {
        browser: true,
        node: true,
    },
    parser: 'vue-eslint-parser',
    parserOptions: {
        parser: '@typescript-eslint/parser',
    },
    extends: [
        '@nuxt/eslint-config',
        'plugin:vue/vue3-recommended',
        'plugin:prettier/recommended',
        'plugin:tailwindcss/recommended',
    ],
    plugins: [],
    rules: {
        'no-console': 0,
        'require-await': 0,
        'no-restricted-syntax': ['error', 'IfStatement > ExpressionStatement > AssignmentExpression'],
        '@typescript-eslint/no-unused-vars': [
            'warn', // or "error"
            {
                argsIgnorePattern: '^_',
                varsIgnorePattern: '^_',
                caughtErrorsIgnorePattern: '^_',
            },
        ],
    },
};

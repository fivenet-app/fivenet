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
    extends: ['@nuxtjs/eslint-config-typescript', 'plugin:vue/vue3-recommended', 'plugin:prettier/recommended'],
    plugins: [],
    rules: {
        'no-console': 0,
        'require-await': 0,
        // Vue
        'vue/multi-word-component-names': 0,
        // vee-validate requires v-bind for custom fields
        'vue/no-deprecated-v-bind-sync': 0,
    },
};

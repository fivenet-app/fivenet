export default defineI18nConfig(nuxt => ({
    legacy: false,
    locales: [
        {
            code: 'en',
            file: 'en-US.json',
        },
        {
            code: 'de',
            file: 'de-DE.json',
        },
    ],
    lazy: true,
    langDir: 'lang',
    defaultLocale: 'en',
}));

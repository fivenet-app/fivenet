export default defineAppConfig({
    ui: {
        primary: 'green',
        gray: 'slate',
        footer: {
            bottom: {
                left: 'text-sm text-gray-500 dark:text-gray-400',
                wrapper: 'border-t border-gray-200 dark:border-gray-800',
            },
        },
    },
    seo: {
        siteName: 'FiveNet',
    },
    toc: {
        bottom: {
            edit: 'https://github.com/galexrt/fivenet/edit/main/content',
        },
    },
    links: {
        imprint: '',
        privacyPolicy: '',
    },
});

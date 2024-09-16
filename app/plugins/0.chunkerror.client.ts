// Modified version based on https://github.com/nuxt/nuxt/pull/19086#issuecomment-1553385289
export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.hook('app:chunkError', ({ error }) => {
        const error_list = [
            'error loading dynamically imported module',
            'Importing a module script failed',
            'Failed to fetch dynamically imported module',
        ];

        if (typeof error.message === 'string')
            for (const message of error_list) {
                if (message.indexOf(error.message) > -1) {
                    reloadNuxtApp({ persistState: true });
                }
            }
    });
});

export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.hook('app:chunkError', ({ error }) => {
        const errorMessages = [
            'error loading dynamically imported module',
            'importing a module script failed',
            'failed to fetch dynamically imported module',
        ];

        const message = String(error.message).toLowerCase();

        if (errorMessages.some((errMsg) => message.includes(errMsg))) {
            reloadNuxtApp({ persistState: true });
            return;
        }
    });
});

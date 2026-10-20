import { recordClientError } from '~/utils/client-errors';

type ErrorCollectionRuntime = { cleanup: () => void };
type WindowWithErrorCollectionRuntime = Window & { __fivenetErrorCollectionRuntime?: ErrorCollectionRuntime };

export default defineNuxtPlugin({
    name: 'error-collection',
    enforce: 'pre',

    setup(nuxtApp) {
        const router = useRouter();
        const windowWithRuntime = window as WindowWithErrorCollectionRuntime;
        let lastErrorKey = '';
        let lastErrorAt = 0;

        const record = (error: unknown, source: string): void => {
            const normalized = error instanceof Error ? `${error.name}:${error.message}:${error.stack ?? ''}` : String(error);
            const now = Date.now();

            // Nuxt and the browser can report the same exception through multiple hooks.
            if (normalized === lastErrorKey && now - lastErrorAt < 1000) return;
            lastErrorKey = normalized;
            lastErrorAt = now;

            recordClientError(error, { source, route: router.currentRoute.value.fullPath });
        };

        nuxtApp.hook('app:error', (error) => record(error, 'Nuxt application'));
        nuxtApp.hook('vue:error', (error, _instance, info) => record(error, `Vue: ${info}`));

        const handleErrorEvent = (event: ErrorEvent): void => record(event.error ?? event.message, 'Browser');
        const handleUnhandledRejection = (event: PromiseRejectionEvent): void => record(event.reason, 'Unhandled promise');
        const cleanup = (): void => {
            window.removeEventListener('error', handleErrorEvent);
            window.removeEventListener('unhandledrejection', handleUnhandledRejection);
            if (windowWithRuntime.__fivenetErrorCollectionRuntime?.cleanup === cleanup) {
                delete windowWithRuntime.__fivenetErrorCollectionRuntime;
            }
        };

        // Reuse the same listeners when Vite HMR re-runs this client plugin.
        windowWithRuntime.__fivenetErrorCollectionRuntime?.cleanup();
        window.addEventListener('error', handleErrorEvent);
        window.addEventListener('unhandledrejection', handleUnhandledRejection);
        windowWithRuntime.__fivenetErrorCollectionRuntime = { cleanup };

        if (import.meta.hot) import.meta.hot.dispose(cleanup);
    },
});

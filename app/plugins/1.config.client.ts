import { defineNuxtPlugin } from '#app';
import type { ServerAppConfig } from '~/typings';

async function loadConfig(): Promise<void> {
    // 7.5 seconds should be enough to retrieve the config from the server...
    const abort = new AbortController();
    const tId = setTimeout(() => abort.abort(), 7.5 * 1000);

    try {
        const resp = await $fetch<ServerAppConfig>('/api/config', {
            method: 'POST',
            signal: abort.signal,
        });

        updateAppConfig(resp);
    } catch (e) {
        const err = e as Error;
        throw createError({
            statusCode: 500,
            statusMessage: 'Failed to get FiveNet config from backend',
            message: err.message + '(Cause: ' + err.cause + ')',
            fatal: true,
            unhandled: false,
        });
    } finally {
        clearTimeout(tId);
    }
}

export default defineNuxtPlugin({
    name: 'config',
    parallel: true,
    async setup(_) {
        await loadConfig();

        return {};
    },
});

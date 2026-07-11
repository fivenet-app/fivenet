import { defineNuxtPlugin } from '#app';
import type { ClientConfig } from '~~/gen/ts/resources/clientconfig/clientconfig';

const appConfigPromise = loadConfig();

function applyClientConfig(clientConfig: ClientConfig): void {
    updateAppConfig({ ...clientConfig });
}

async function loadConfig(): Promise<ClientConfig> {
    // If running in Vitest, return empty config
    if (import.meta.env.VITEST) return {} as ClientConfig;

    const abort = new AbortController();
    const tId = setTimeout(() => abort.abort(), 7_500);

    try {
        return await $fetch<ClientConfig>('/api/config', {
            method: 'POST',
            signal: abort.signal,
        });
    } catch (e) {
        console.error('Failed to get FiveNet config from backend', e);
        const err = e as Error;
        throw createError({
            status: 500,
            statusText: 'Failed to get FiveNet config from backend',
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
    parallel: false,
    enforce: 'post',

    async setup() {
        applyClientConfig(await appConfigPromise);

        return {
            provide: {
                appConfigPromise: appConfigPromise,
                reloadConfigFromServer: loadConfig,
                applyClientConfig,
            },
        };
    },
});

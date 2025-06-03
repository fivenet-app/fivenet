import { defineNuxtPlugin } from '#app';
import type { ServerAppConfig } from '~/typings';

const appConfigPromise = loadConfig();

async function loadConfig(): Promise<ServerAppConfig> {
    const abort = new AbortController();
    const tId = setTimeout(() => abort.abort(), 7_500);

    try {
        const resp = await $fetch<ServerAppConfig>('/api/config', {
            method: 'POST',
            signal: abort.signal,
        });
        updateAppConfig(resp);
        return resp;
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
    async setup(nuxtApp) {
        await appConfigPromise;

        nuxtApp.provide('appConfigPromise', appConfigPromise);
    },
});

import { defineNuxtPlugin } from '#app';
import type { ClientConfig } from '~~/gen/ts/resources/clientconfig/clientconfig';

const appConfigPromise = loadConfig();

async function loadConfig(): Promise<ClientConfig> {
    const abort = new AbortController();
    const tId = setTimeout(() => abort.abort(), 7_500);

    try {
        const resp = await $fetch<ClientConfig>('/api/config', {
            method: 'POST',
            signal: abort.signal,
        });
        updateAppConfig({ ...resp });
        return resp;
    } catch (e) {
        console.error('Failed to get FiveNet config from backend', e);
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

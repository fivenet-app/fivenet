import { defineNuxtPlugin } from '#app';
import type { AppConfig } from 'nuxt/schema';

async function loadConfig(): Promise<void> {
    // 7.5 seconds should be enough to retrieve the config
    const abort = new AbortController();
    const tId = setTimeout(() => abort.abort(), 7.5 * 1000);

    const resp = await $fetch<AppConfig>('/api/config', {
        method: 'POST',
        signal: abort.signal,
    })
        .catch((e) => {
            throw createError({
                statusCode: 500,
                statusMessage: 'Failed to get FiveNet config from backend',
                message: e,
                fatal: true,
                unhandled: false,
            });
        })
        .finally(() => clearTimeout(tId));

    updateAppConfig(resp);
}

export default defineNuxtPlugin({
    name: 'config',
    parallel: true,
    async setup(_) {
        await loadConfig();
        return {};
    },
});

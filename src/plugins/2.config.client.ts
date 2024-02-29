import { defineNuxtPlugin } from '#app';
import type { AppConfig } from '~/shims';

async function loadConfig(): Promise<void> {
    // 10 seconds should be enough to retrieve the config
    const abort = new AbortController();
    const tId = setTimeout(() => abort.abort(), 10 * 1000);

    const resp = await fetch('/api/config', {
        method: 'POST',
        signal: abort.signal,
    });
    clearTimeout(tId);

    if (!resp.ok) {
        const text = await resp.text();
        throw createError({
            statusCode: 500,
            statusMessage: 'Failed to get FiveNet config from backend',
            message: text,
            fatal: true,
            unhandled: false,
        });
    }
    const data = (await resp.json()) as AppConfig;
    updateAppConfig(data);
}

export default defineNuxtPlugin({
    async setup(_) {
        await loadConfig();
    },
});

import { defineNuxtPlugin } from '#app';
import { type NuxtError } from 'nuxt/app';
import type { AppConfig } from '~/shims';

async function loadConfig(): Promise<void> {
    try {
        // 6 seconds should be enough
        const abort = new AbortController();
        const tId = setTimeout(() => abort.abort(), 8000);

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
    } catch (e) {
        showError(e as NuxtError);
        throw e;
    }
}

export default defineNuxtPlugin({
    async setup(_) {
        await loadConfig();
    },
});

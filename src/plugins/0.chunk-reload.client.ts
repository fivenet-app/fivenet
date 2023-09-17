// Taken from https://github.com/nuxt/nuxt/blob/main/packages/nuxt/src/app/plugins/chunk-reload.client.ts

import { reloadNuxtApp } from '#app/composables/chunk';
import { useRouter } from '#app/composables/router';
import { defineNuxtPlugin, useRuntimeConfig } from '#app/nuxt';
import { joinURL } from 'ufo';

export default defineNuxtPlugin({
    name: 'nuxt:chunk-reload',
    setup(nuxtApp) {
        const router = useRouter();
        const config = useRuntimeConfig();

        const chunkErrors = new Set();

        router.beforeEach(() => {
            chunkErrors.clear();
        });
        nuxtApp.hook('app:chunkError', ({ error }) => {
            chunkErrors.add(error);
        });

        router.onError((error, to) => {
            if (chunkErrors.has(error)) {
                const isHash = 'href' in to && (to.href as string).startsWith('#');
                const path = isHash ? config.app.baseURL + (to as any).href : joinURL(config.app.baseURL, to.fullPath);
                reloadNuxtApp({ path, persistState: true });
            }
        });
    },
});

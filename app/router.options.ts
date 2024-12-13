// Based upon Nuxt UI Pro's router options to fix the "not a valid query selector" DOM Exceptions

import type { RouterConfig } from '@nuxt/schema';

function findHashPosition(hash: string): undefined {
    try {
        const el = document.querySelector(hash);
        // vue-router does not incorporate scroll-margin-top on its own.
        if (el) {
            el.scrollIntoView({
                behavior: 'smooth',
            });
        }
    } catch (e) {
        console.warn('Query Selector exception', e);
    }
}

// https://router.vuejs.org/api/#routeroptions
export default <RouterConfig>{
    scrollBehavior(to, from, savedPosition) {
        const nuxtApp = useNuxtApp();

        if (history.state && history.state.stop) {
            return;
        }
        if (history.state && history.state.smooth) {
            return {
                el: history.state.smooth,
                behavior: 'smooth',
            };
        }

        // If history back
        if (savedPosition) {
            // Handle Suspense resolution
            return new Promise((resolve) => {
                nuxtApp.hooks.hookOnce('page:finish', () => {
                    setTimeout(() => resolve(savedPosition), 50);
                });
            });
        }

        // Scroll to heading on click
        if (to.hash && to.hash !== '#') {
            return new Promise((resolve) => {
                if (to.path === from.path) {
                    setTimeout(() => resolve(findHashPosition(to.hash)), 50);
                } else {
                    nuxtApp.hooks.hookOnce('page:finish', () => {
                        setTimeout(() => resolve(findHashPosition(to.hash)), 50);
                    });
                }
            });
        }

        // Scroll to top of window after page loaded
        return new Promise((resolve) => {
            nuxtApp.hooks.hookOnce('page:finish', () => {
                setTimeout(() => resolve({ top: 0 }), 50);
            });
        });
    },
};
